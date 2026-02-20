package main

import (
	"fmt"
	"minimal/minimal-core/built-in/extensions/commands/templates"
	"minimal/minimal-core/built-in/extensions/stores/directory"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	commands.AddCommand("hello", helloWorld)

	projectCreator := templates.NewProjectCreator(sourceGen)
	projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(projectCreator.SourceGen))
	commands.AddCommand("new", projectCreator.NewProject)
}

func helloWorld(_ []string) bool {
	fmt.Println("Hello, world!")
	return true
}
