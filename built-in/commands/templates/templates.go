package templates

import (
	"errors"
	"flag"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"

	"github.com/google/shlex"
	"github.com/rs/zerolog"
)

const (
	destinationFlagName = "destination"
	implementationArgsFlagName = "args"
)

var DuplicateTemplateStore = errors.New("template store with this name already exists")

type ProjectCreator struct {
	logger          zerolog.Logger
	SourceGen       *logging.SourceGenerator
	stores          map[string]TemplateStore
}

type TemplateStore interface {
	Name() string
	HasTemplate(name string) bool
	LoadTemplate(name, projectName, destination string, args []string) (ok bool)
}

type MutableTemplateStore interface {
	TemplateStore
	StoreTemplate(name, source string, args []string) (ok bool)
}

func NewProjectCreator(sourceGen *logging.SourceGenerator) *ProjectCreator {
	logger, gen := sourceGen.GetLogger("templates")
	
	return &ProjectCreator{logger, gen, make(map[string]TemplateStore, 0)}
}

func (p *ProjectCreator) RegisterTemplateStore(store TemplateStore) error {
	if _, ok := p.stores[store.Name()]; ok {
		p.logDuplicateTemplateStore(store.Name())
		return DuplicateTemplateStore
	}
	
	p.stores[store.Name()] = store
	p.logTemplateRegistered(store.Name())

	return nil
}

func (p *ProjectCreator) NewProject(args []string) bool {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // TODO log error

	var destination string
	fs.StringVar(&destination, destinationFlagName, "", "")
	fs.StringVar(&destination, string(destinationFlagName[0]), "", "")

	var implementationArgs string
	fs.StringVar(&implementationArgs, implementationArgsFlagName, "", "")
	fs.StringVar(&implementationArgs, string(implementationArgsFlagName[0]), "", "")

	if err := fs.Parse(args); err != nil {
        return false
    }

	implArgs, err := shlex.Split(implementationArgs)

	if err != nil {
		// TODO log error
	}

	var templateName string
	var projectName string

	switch fs.NArg() {
	case 1:
		projectName = fs.Arg(0)
	case 2:
		templateName = fs.Arg(0)
		projectName = fs.Arg(1)
	default:
		p.logger.Error().
			Int("expected_args_1", 1).
			Int("expected_args_2", 2).
			Int("actual_args", fs.NArg()).
			Msg("incorrect amount of arguments")

		return false
	}

	availableSources := make([]TemplateStore, 0)

	for _, store := range p.stores {
		if store.HasTemplate(templateName) {
			availableSources = append(availableSources, store)
		}
	}

	switch len(availableSources) {
	case 0:
		p.logger.Error().Str("template_name", templateName).Msg("no sources to load template")
		return false
	case 1:
		if ok := availableSources[0].LoadTemplate(templateName, projectName, destination, implArgs); !ok {
			p.logger.Error().Str("template_name", templateName).Str("source", availableSources[0].Name()).Msg("failure during template load")
		}
	default:
		// TODO ask user which store to load from
	}

	return true
}

func (p *ProjectCreator) StoreTemplate(args []string) bool {
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

	availableStores := make([]MutableTemplateStore, 0)

	for _, store := range p.stores {
		if mutableStore, ok := store.(MutableTemplateStore); ok {
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
		Err(DuplicateTemplateStore).
		Msg("")
}

func (p *ProjectCreator) logTemplateRegistered(name string) {
	p.logger.Debug().
		Str("template_name", name).
		Msg("template registered")
}
