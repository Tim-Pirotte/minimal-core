package setup

import (
	"minimal/minimal-core/built-in/commands/templates"
	messaging "minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/stores/directory"
	"minimal/minimal-core/built-in/ui"
	"minimal/minimal-core/built-in/user-interfaces/tui"
	"os"
)

func RegisterCommands(
    commands *startup.Commands,
    messenger *messaging.Messenger,
    ui ui.UI,
) {
    projectCreator := templates.NewProjectCreator(messenger, ui)
    err := projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(messenger), "directory", 1)

    if err != nil {
        messenger.Send(messaging.Message{
            Message: "Failed to register the directory template store",
            Severity: messaging.Error,
        })

        os.Exit(1)
    }

    err = commands.AddCommand("new", projectCreator.NewProjectCLI)

    if err != nil {
        messenger.Send(messaging.Message{
            Message: "Failed to register the 'new' command",
            Severity: messaging.Error,
        })

        os.Exit(1)
    }

    err = commands.AddCommand("tui", tui.NewTUI().StartTUI)

    if err != nil {
        messenger.Send(messaging.Message{Message: "Failed to register the 'tui' command"})

        os.Exit(1)
    }
}
