package main

import (
	"fmt"
	"minimal/minimal-core/built-in/commands/templates"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	commands.AddCommand("hello", helloWorld)

	projectCreator := templates.NewProjectCreator(sourceGen)
	commands.AddCommand("new", projectCreator.NewProject)
}

func helloWorld() {
	fmt.Println("Hello, world!")
}
