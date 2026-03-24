package templates

import (
	"errors"
	"fmt"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/ui"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"strconv"
	"strings"
)

var ErrDuplicateTemplateStore = errors.New("template store with this name and priority already exists")

type ProjectCreator struct {
	logger          logging.Logger
	messenger       *usermessaging.Messenger
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
	StoreTemplate(name, source string, fields map[string]string) (ok bool)
	RemoveTemplate(name string, args []string) (ok bool)
}

func NewProjectCreator(sourceGen *logging.SourceGenerator, messenger *usermessaging.Messenger, ui ui.UI) *ProjectCreator {
	logger, _ := sourceGen.GetLogger("templates")
	
	return &ProjectCreator{logger, messenger, ui, make([]templateStore, 0)}
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
	newProjectArgs := p.getNewProjectArgs(args)

	var templateName string
	var projectName string

	numberOfArguments := uint(len(args)) - newProjectArgs.positionalArgStart

	switch numberOfArguments {
	case 1:
		projectName = args[newProjectArgs.positionalArgStart]
	case 2:
		templateName = args[newProjectArgs.positionalArgStart]
		projectName = args[newProjectArgs.positionalArgStart + 1]
	default:
		p.logIncorrectArguments(int(numberOfArguments))

		return false
	}

	return p.CreateNewProject(templateName, projectName, newProjectArgs.destination, newProjectArgs.fields)
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
				p.logMissingArgumentAfterDestination(args)
			}
		case "-field", "-f":
			if i + 2 < len(args) {
				key := args[i+1]
                val := args[i+2]
                
				fields[key] = val

				i += 2

				positionalArgStart = uint(i) + 1
			} else {
				p.logMissingArgumentsAfterField(args)
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
		p.logNoSourcesForLoading(templateName)
		
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
		p.logger.Error().Str("template_name", templateName).Str("store_name", selectedStore.name).Msg("failure during template load")
		
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
	loadProjectArgs := p.getNewProjectArgs(args)

	var source string
	var templateName string

	numberOfArguments := uint(len(args)) - loadProjectArgs.positionalArgStart

	switch numberOfArguments {
	case 1:
		templateName = args[loadProjectArgs.positionalArgStart]
	case 2:
		source = args[loadProjectArgs.positionalArgStart]
		templateName = args[loadProjectArgs.positionalArgStart + 1]
	default:
		p.logger.Error().
			Int("expected_args_1", 1).
			Int("expected_args_2", 2).
			Int("actual_args", int(numberOfArguments)).
			Msg("incorrect amount of arguments")

		return false
	}

	if ok := p.StoreTemplate(templateName, source, loadProjectArgs.fields); !ok {
		p.logger.Error().Str("template_name", templateName).Str("source_location", source).Msg("failure during template store")
		
		return false
	}

	return true
}

type storeProjectArgs struct {
	fields map[string]string
	positionalArgStart uint
}

func (p *ProjectCreator) getStoreProjectArgs(args []string) storeProjectArgs {
	fields := map[string]string{}
	positionalArgStart := uint(0)
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		switch arg {
		case "-field", "-f":
			if i + 2 < len(args) {
				key := args[i+1]
                val := args[i+2]
                
				fields[key] = val

				i += 2

				positionalArgStart = uint(i) + 1
			} else {
				p.logMissingArgumentsAfterField(args)
			}
		}
	}

	return storeProjectArgs{fields, positionalArgStart}
}

func (p *ProjectCreator) StoreTemplate(templateName, source string, fields map[string]string) bool {
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
		availableStores[0].StoreTemplate(templateName, source, fields)
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

func (p *ProjectCreator) logIncorrectArguments(actual int) {
	p.logger.Error().
		Int("expected_args_1", 1).
		Int("expected_args_2", 2).
		Int("actual_args", actual).
		Msg("incorrect amount of arguments")
	
	t := p.messenger.CreateLogTransaction()
	
	p.messenger.LogMessage(
		t, 
		usermessaging.Message{
			Severity: usermessaging.Critical, 
			Category: "TemplateError", 
			Message: fmt.Sprintf("Expected 1 or 2 arguments but got %d arguments", actual),
		})

	p.messageCLIHelpText(t)
	p.messenger.CommitLogTransaction(t)
}

func (p *ProjectCreator) logMissingArgumentAfterDestination(args []string) {
	p.logger.Error().Strs("args", args).Msg("missing argument after destination flag")

	t := p.messenger.CreateLogTransaction()

	p.messenger.LogMessage(
		t, 
		usermessaging.Message{
			Severity: usermessaging.Error, 
			Category: "TemplateError", 
			Message: "Missing argument after the destination flag",
		},
	)

	p.messageCLIHelpText(t)
	p.messenger.CommitLogTransaction(t)
}

func (p *ProjectCreator) logMissingArgumentsAfterField(args []string) {
	p.logger.Error().Strs("args", args).Msg("missing key and/or value after field flag")

	t := p.messenger.CreateLogTransaction()

	p.messenger.LogMessage(
		t, 
		usermessaging.Message{
			Severity: usermessaging.Error, 
			Category: "TemplateError", 
			Message: "Missing key and/or value after field flag",
		},
	)

	p.messenger.CommitLogTransaction(t)
}

func (p *ProjectCreator) logNoSourcesForLoading(templateName string) {
	p.logger.Error().Str("template_name", templateName).Msg("no sources to load template")

	t := p.messenger.CreateLogTransaction()

	p.messenger.LogMessage(
		t, 
		usermessaging.Message{
			Severity: usermessaging.Critical, 
			Category: "TemplateError", 
			Message: "No available sources to load the template from",
		},
	)

	var hint strings.Builder

	hint.WriteString("These are all the template stores: ")

	for i, s := range p.stores {
		if i != 0 {
			hint.WriteString(", ")
		}
		
		hint.WriteString(s.name)
	}

	p.messenger.LogHint(
		t, 
		usermessaging.Hint{
			Text: hint.String(), 
			MoreInfoReference: "",
		},
	)

	p.messenger.CommitLogTransaction(t)
}

func (p *ProjectCreator) logMissingArgumentAfterSource(args []string) {
	p.logger.Error().Strs("args", args).Msg("missing argument after source flag")

	t := p.messenger.CreateLogTransaction()

	p.messenger.LogMessage(
		t, 
		usermessaging.Message{
			Severity: usermessaging.Error, 
			Category: "TemplateError", 
			Message: "Missing argument after the source flag",
		},
	)

	p.messageCLIHelpText(t)
	p.messenger.CommitLogTransaction(t)
}

func (p *ProjectCreator) messageCLIHelpText(t *usermessaging.Transaction) {
	// TODO move the explaination of the complete command to the info reference.
	p.messenger.LogHint(
		t, 
		usermessaging.Hint{
			Text: 
				"'mnm new' followed by a project name or a template name and a project name. " +
				"A destination can be set using -d or -destination followed by a path." +
				"Fields passed to template stores can be set using -f or -field followed by a key and a value.", 
			MoreInfoReference: "",
		},
	)
}
