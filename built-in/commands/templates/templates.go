package templates

import (
	"errors"
	"flag"
	logging "minimal/minimal-core/built-in/internal-logging"

	"github.com/rs/zerolog"
)

const destinationFlagName = "destination"

var DuplicateTemplateStore = errors.New("template store with this name already exists")

type ProjectCreator struct {
	logger    zerolog.Logger
	SourceGen *logging.SourceGenerator
	stores    map[string]TemplateStore
}

type TemplateStore interface {
	HasTemplate(name string) bool
	LoadTemplate(name, projectName, destination string) (ok bool)
	StoreTemplate(name, location string) (ok bool)
}

func NewProjectCreator(sourceGen *logging.SourceGenerator) *ProjectCreator {
	logger, gen := sourceGen.GetLogger("templates")
	
	return &ProjectCreator{logger, gen, make(map[string]TemplateStore, 0)}
}

func (p *ProjectCreator) AddTemplateStore(store TemplateStore, name string) error {
	if _, ok := p.stores[name]; ok {
		p.logDuplicateTemplateStore(name)
		return DuplicateTemplateStore
	}
	
	p.stores[name] = store
	p.logTemplateRegistered(name)

	return nil
}

func (p *ProjectCreator) NewProject() {
	var destination string
	flag.StringVar(&destination, destinationFlagName, "", "")
	flag.StringVar(&destination, string(destinationFlagName[0]), "", "")

	flag.Parse()

	var templateName string

	switch flag.NArg() {
	case 1:
	case 2:
		templateName = flag.Arg(0)
	default:
		p.logger.Error().
			Int("expected_args_1", 1).
			Int("expected_args_2", 2).
			Int("actual_args", flag.NArg()).
			Msg("incorrect amount of arguments")

		return
	}

	availableSources := make([]TemplateStore, 0)

	for _, store := range p.stores {
		if store.HasTemplate(templateName) {
			availableSources = append(availableSources, store)
		}
	}

	switch len(availableSources) {
	case 0:
		// TODO log error
	case 1:
		if ok := availableSources[0].LoadTemplate(templateName, "", destination); !ok {
			// TODO log error
		}
	default:
		// TODO ask user which store to load from
	}
}

func CreateTemplate() {
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
