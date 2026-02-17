package main

import (
	"fmt"
	"minimal/minimal-core/built-in/commands/templates"
	"minimal/minimal-core/built-in/extensions/stores/directory"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	commands.AddCommand("hello", helloWorld)

	projectCreator := templates.NewProjectCreator(sourceGen)
	projectCreator.AddTemplateStore(directory.NewDirectoryStore(projectCreator.SourceGen))
	commands.AddCommand("new", projectCreator.NewProject)
}

func helloWorld() bool {
	fmt.Println("Hello, world!")
	return true
}
