package templates

import (
	"errors"
	"flag"
	"fmt"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"

	"github.com/rs/zerolog"
)

const destinationFlagName = "destination"

var DuplicateTemplateStore = errors.New("template store with this name already exists")

type ProjectCreator struct {
	logger          zerolog.Logger
	SourceGen       *logging.SourceGenerator
	stores          map[string]TemplateStore
}

type TemplateStore interface {
	Name() string
	HasTemplate(name string) bool
	LoadTemplate(name, projectName, destination string) (ok bool)
}

type MutableTemplateStore interface {
	TemplateStore
	StoreTemplate(name, location string) (ok bool)
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

func (p *ProjectCreator) NewProject() bool {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)

	var destination string
	fs.StringVar(&destination, destinationFlagName, "", "")
	fs.StringVar(&destination, string(destinationFlagName[0]), "", "")

	fmt.Println(os.Args, len(os.Args))

	if err := fs.Parse(os.Args[2:]); err != nil {
        return false
    }

	var templateName string
	var projectName string

	fmt.Println(fs.NArg())

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
		if ok := availableSources[0].LoadTemplate(templateName, projectName, destination); !ok {
			p.logger.Error().Str("template_name", templateName).Str("source", availableSources[0].Name()).Msg("failure during template load")
		}
	default:
		// TODO ask user which store to load from
	}

	return true
}

func StoreTemplate() {
	var symbolicLink bool
	flag.BoolVar(&symbolicLink, "ln", false, "")

	flag.Parse()
}

func saveTemplate() {

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
