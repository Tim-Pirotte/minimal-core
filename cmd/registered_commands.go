package main

import (
	"fmt"
	"minimal/minimal-core/built-in/extensions/commands/templates"
	logrendering "minimal/minimal-core/built-in/extensions/outputs/log-renderer"
	"minimal/minimal-core/built-in/extensions/stores/directory"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"
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

	err = commands.AddCommand("message", messagingTest)

	if err != nil {
		// TODO log error
	}
}

func helloWorld(_ []string) bool {
	fmt.Println("Hello, world!")
	return true
}

func messagingTest(_ []string) bool {
	logRenderer := logrendering.NewLogger(os.Stdout, logrendering.Config{})
	messaging := usermessaging.NewMessenger()
	defer messaging.Close()
	messaging.AddOutput(logRenderer)

	t := messaging.CreateLogTransaction()
	messaging.LogMessage(t, usermessaging.Message{
		Severity: usermessaging.Critical,
		Category: "Parsing error", 
		Message: "Hello, World!",
	})
	messaging.LogHint(t, usermessaging.Hint{Text: "Test", MoreInfoReference: "More info test"})
	messaging.CommitLogTransaction(t)

	return true
}
