package setup

import (
	"minimal/minimal-lang/built-in/commands/templates"
	messaging "minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/startup"
	"minimal/minimal-lang/built-in/stores/directory"
	"minimal/minimal-lang/built-in/ui"
	"minimal/minimal-lang/built-in/user-interfaces/tui"
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
