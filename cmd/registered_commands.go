package main

import (
	"minimal/minimal-core/built-in/commands/templates"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/stores/directory"
	"minimal/minimal-core/built-in/user-interfaces/tui"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	logger, _ := sourceGen.GetLogger("commandRegistry")

	projectCreator := templates.NewProjectCreator(sourceGen)
	err := projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(projectCreator.SourceGen))

	if err != nil {
		logger.Fatal().Msg("Failed to register the directory template store")
	}

	err = commands.AddCommand("new", projectCreator.NewProject)

	if err != nil {
		logger.Fatal().Msg("Failed to register the new command")
	}

	err = commands.AddCommand("tui", tui.NewTUI().StartTUI)

	if err != nil {
		logger.Fatal().Msg("Failed to register the tui command")
	}
}
