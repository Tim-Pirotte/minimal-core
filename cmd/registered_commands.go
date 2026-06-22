package main

import (
	"minimal/minimal-core/built-in/commands/templates"
	logging "minimal/minimal-core/built-in/internal-logging"
	messaging "minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/stores/directory"
	"minimal/minimal-core/built-in/ui"
	"minimal/minimal-core/built-in/user-interfaces/tui"
)

func registerCommands(
	commands *startup.Commands,
	sourceGen *logging.SourceGenerator,
	messenger *messaging.Messenger,
	ui ui.UI,
) {
	logger, _ := sourceGen.GetLogger("CommandRegistry")

	projectCreator := templates.NewProjectCreator(sourceGen, messenger, ui)
	err := projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(sourceGen), "directory", 1)

	if err != nil {
		logger.Fatal().Msg("Failed to register the directory template store")
	}

	err = commands.AddCommand("new", projectCreator.NewProjectCLI)

	if err != nil {
		logger.Fatal().Msg("Failed to register the new command")
	}

	err = commands.AddCommand("tui", tui.NewTUI().StartTUI)

	if err != nil {
		logger.Fatal().Msg("Failed to register the tui command")
	}
}
