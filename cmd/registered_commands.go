package main

import (
	"fmt"
	"minimal/minimal-core/built-in/extensions/commands/templates"
	"minimal/minimal-core/built-in/extensions/stores/directory"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	err := commands.AddCommand("hello", helloWorld)

	if err != nil {
		// TODO log error
	}

	projectCreator := templates.NewProjectCreator(sourceGen)
	err = projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(projectCreator.SourceGen))

	if err != nil {
		// TODO log error
	}

	err = commands.AddCommand("new", projectCreator.NewProject)

	if err != nil {
		// TODO log error
	}
}

func helloWorld(_ []string) bool {
	fmt.Println("Hello, world!")
	return true
}
