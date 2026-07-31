package templates

import (
	"errors"
	messaging "minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/ui"
	"strconv"
	"strings"
)

var ErrDuplicateTemplateStore = errors.New("template store with this name and priority already exists")

type ProjectCreator struct {
    messenger       *messaging.Messenger
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

func NewProjectCreator(messenger *messaging.Messenger, ui ui.UI) *ProjectCreator {
    return &ProjectCreator{messenger, ui, make([]templateStore, 0)}
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
        p.logIncorrectLoadArguments(int(numberOfArguments))

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
        p.messenger.Send(messaging.Message{
            Message: "Could not load the template '" + templateName + "' from '" + selectedStore.name + "'",
            Severity: messaging.Error,
        })

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
        p.messenger.Send(messaging.Message{
            Message: "No valid answer for project creation",
            Severity: messaging.Error,
        })

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
    loadProjectArgs := p.getStoreProjectArgs(args)

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
        p.logIncorrectStoreArguments(int(numberOfArguments))

        return false
    }

    if ok := p.StoreTemplate(templateName, source, loadProjectArgs.fields); !ok {
        p.messenger.Send(messaging.Message{
            Message: "Storing the template failed",
            Severity: messaging.Error,
        })

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
    availableStores := make([]MutableTemplateStore, 0)

    for _, s := range p.stores {
        var store TemplateStore = s

        if mutableStore, ok := store.(MutableTemplateStore); ok {
            availableStores = append(availableStores, mutableStore)
        }
    }

    switch len(availableStores) {
    case 0:
        p.logNoSourcesForStoring()

        return false
    case 1:
        availableStores[0].StoreTemplate(templateName, source, fields)
    default:
        // TODO ask user where to store
    }

    return true
}

func (p *ProjectCreator) logDuplicateTemplateStore(name string) {
    p.messenger.Send(messaging.Message{
        Message: "A template with the name '" + name + "' already exists in this store",
        Severity: messaging.Error,
    })
}

func (p *ProjectCreator) logTemplateRegistered(name string) {
    p.messenger.Send(messaging.Message{
        Message: "A template with the name '" + name + "' has successfully been registered",
        Severity: messaging.Info,
    })
}

func (p *ProjectCreator) logIncorrectLoadArguments(actual int) {
    p.messenger.Send(messaging.Message{
        Message: "Expected 1 or 2 arguments but got " + strconv.Itoa(actual),
        Severity: messaging.Error,
    })
}

func (p *ProjectCreator) logMissingArgumentAfterDestination(args []string) {
    p.messenger.Send(messaging.Message{
        Message: "Missing argument after destination flag",
        Severity: messaging.Error,
    })
}

func (p *ProjectCreator) logMissingArgumentsAfterField(args []string) {
    p.messenger.Send(messaging.Message{
        Message: "Missing key after field flag",
        Severity: messaging.Error,
    })
}

func (p *ProjectCreator) logNoSourcesForLoading(templateName string) {
    var hint strings.Builder

    hint.WriteString("These are all the template stores: ")

    for i, s := range p.stores {
        if i != 0 {
            hint.WriteString(", ")
        }

        hint.WriteString(s.name)
    }

    p.messenger.Send(messaging.Message{
        Message: "The template '" + templateName + "' could not be found",
        Severity: messaging.Error,
        Notes: []string{hint.String()},
    })
}

func (p *ProjectCreator) logIncorrectStoreArguments(actual int) {
    p.messenger.Send(messaging.Message{
        Message: "Expected 1 or 2 arguments but got " + strconv.Itoa(actual),
        Severity: messaging.Error,
    })
}

func (p *ProjectCreator) logNoSourcesForStoring() {
    p.messenger.Send(messaging.Message{
        Message: "There are no stores to save the template in",
        Severity: messaging.Error,
    })
}
