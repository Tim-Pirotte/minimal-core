package templates

import (
	"errors"
	"flag"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/ui"
	"strconv"
	"strings"

	"github.com/google/shlex"
)

const (
	destinationFlagName = "destination"
	implementationArgsFlagName = "args"
)

var ErrDuplicateTemplateStore = errors.New("template store with this name and priority already exists")

type ProjectCreator struct {
	logger          logging.Logger
	ui              ui.UI
	stores          []templateStore
}

type templateStore struct {
	TemplateStore
	name string
	priority uint
}

type TemplateStore interface {
	HasTemplate(name string) bool
	LoadTemplate(name, projectName, destination string, fields map[string]string) (ok bool)
}

type MutableTemplateStore interface {
	TemplateStore
	StoreTemplate(name, source string, args []string) (ok bool)
	RemoveTemplate(name string, args []string) (ok bool)
}

func NewProjectCreator(sourceGen *logging.SourceGenerator, ui ui.UI) *ProjectCreator {
	logger, _ := sourceGen.GetLogger("templates")
	
	return &ProjectCreator{logger, ui, make([]templateStore, 0)}
}

func (p *ProjectCreator) RegisterTemplateStore(store TemplateStore, name string, priority uint) error {
	for _, s := range p.stores {
		if s.name == name && s.priority == priority {
			p.logDuplicateTemplateStore(name)

			return ErrDuplicateTemplateStore
		}
	}
	
	p.stores = append(p.stores, templateStore{store, name, priority})
	p.logTemplateRegistered(name)

	return nil
}

func (p *ProjectCreator) NewProjectCLI(args []string) bool {
	projectArgs := p.getNewProjectArgs(args)

	var templateName string
	var projectName string

	switch uint(len(args)) - projectArgs.positionalArgStart {
	case 1:
		projectName = args[projectArgs.positionalArgStart]
	case 2:
		templateName = args[projectArgs.positionalArgStart]
		projectName = args[projectArgs.positionalArgStart + 1]
	default:
		p.logger.Error().
			Int("expected_args_1", 1).
			Int("expected_args_2", 2).
			Int("actual_args", len(args) - int(projectArgs.positionalArgStart)).
			Msg("incorrect amount of arguments")

		return false
	}

	return p.CreateNewProject(templateName, projectName, projectArgs.destination, projectArgs.fields)
}

type newProjectArgs struct {
	destination string
	fields map[string]string
	positionalArgStart uint
}

func (p *ProjectCreator) getNewProjectArgs(args []string) newProjectArgs {
	destination := ""
	fields := map[string]string{}
	positionalArgStart := uint(0)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		switch arg {
		case "-destination", "-d":
			if i + 1 < len(args) {
				destination = args[i + 1]
				i++

				positionalArgStart = uint(i) + 1
			} else {
				p.logger.Error().Strs("args", args).Msg("missing argument after destination flag")
			}
		case "-field", "-f":
			if i + 2 < len(args) {
				key := args[i+1]
                val := args[i+2]
                
				fields[key] = val

				i += 2

				positionalArgStart = uint(i) + 1
			} else {
				p.logger.Error().Strs("args", args).Msg("missing key and/or value after field flag")
			}
		}
	}

	return newProjectArgs{destination, fields, positionalArgStart}
}

func (p *ProjectCreator) CreateNewProject(templateName, projectName, destination string, fields map[string]string) bool {
	availableSources := p.getStoresWithTemplate(templateName)

	var selectedStore templateStore

	switch len(availableSources) {
	case 0:
		p.logger.Error().Str("template_name", templateName).Msg("no sources to load template")
		
		return false
	case 1:
		selectedStore = availableSources[0]
	default:
		var ok bool

		if selectedStore, ok = p.askUserForStore(availableSources); !ok {
			return false
		}
	}

	if ok := selectedStore.LoadTemplate(templateName, projectName, destination, fields); !ok {
		p.logger.Error().Str("template_name", templateName).Str("source", selectedStore.name).Msg("failure during template load")
		
		return false
	}

	return true
}

func (p *ProjectCreator) getStoresWithTemplate(templateName string) []templateStore {
	highestPriority := uint(0)
	availableSources := make([]templateStore, 0)

	for _, store := range p.stores {
		if store.HasTemplate(templateName) && store.priority >= highestPriority {
			if store.priority > highestPriority {
				availableSources = make([]templateStore, 0)
				highestPriority = store.priority
			}

			availableSources = append(availableSources, store)
		}
	}

	return availableSources
}

func (p *ProjectCreator) askUserForStore(stores []templateStore) (templateStore, bool) {
	var question strings.Builder
	question.WriteString("Which store would you like to load from? ")

	for i, source := range stores {
		question.WriteString("(")
		question.WriteString(strconv.Itoa(i))
		question.WriteString(": ")
		question.WriteString(source.name)
		question.WriteString(")")
	}

	answer, ok := p.ui.PromptString(question.String(), "0")

	if !ok {
		p.logger.Error().Str("answer", answer).Msg("no valid answer for project creation")

		return templateStore{}, false
	}

	answerAsInt, err := strconv.Atoi(answer)

	if err == nil && 0 <= answerAsInt && answerAsInt < len(stores) {
		return stores[answerAsInt], true
	} else {
		for _, s := range stores {
			if s.name == answer {
				return s, true
			}
		}
		
		return templateStore{}, false
	}
}

func (p *ProjectCreator) StoreTemplateCLI(args []string) bool {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // TODO log error

	var implementationArgs string
	fs.StringVar(&implementationArgs, implementationArgsFlagName, "", "")
	fs.StringVar(&implementationArgs, string(implementationArgsFlagName[0]), "", "")

	if err := fs.Parse(args); err != nil {
        return false
    }

	implArgs, err := shlex.Split(implementationArgs)

	if err != nil {
		// TODO log error
		return false
	}

	var source string
	var templateName string

	switch fs.NArg() {
	case 1:
		templateName = fs.Arg(0)
	case 2:
		source = fs.Arg(0)
		templateName = fs.Arg(1)
	default:
		p.logger.Error().
			Int("expected_args_1", 1).
			Int("expected_args_2", 2).
			Int("actual_args", fs.NArg()).
			Msg("incorrect amount of arguments")

		return false
	}

	return p.StoreTemplate(templateName, source, implArgs)
}

func (p *ProjectCreator) StoreTemplate(templateName, source string, implArgs []string) bool {
	highestPriority := uint(0)
	availableStores := make([]MutableTemplateStore, 0)

	for _, s := range p.stores {
		var store TemplateStore = s

		if mutableStore, ok := store.(MutableTemplateStore); ok && s.priority >= highestPriority {
			if s.priority > highestPriority {
				availableStores = make([]MutableTemplateStore, 0)
				highestPriority = s.priority
			}

			availableStores = append(availableStores, mutableStore)
		}
	}
	
	switch len(availableStores) {
	case 0:
		// TODO log error
	case 1:
		availableStores[0].StoreTemplate(templateName, source, implArgs)
	default:
		// TODO ask user where to store
	}

	return true
}

func (p *ProjectCreator) logDuplicateTemplateStore(name string) {
	p.logger.Error().
		Str("template_name", name).
		Err(ErrDuplicateTemplateStore).
		Msg("")
}

func (p *ProjectCreator) logTemplateRegistered(name string) {
	p.logger.Debug().
		Str("template_name", name).
		Msg("template registered")
}
