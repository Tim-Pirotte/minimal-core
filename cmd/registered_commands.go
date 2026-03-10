package main

import (
	"fmt"
	"minimal/minimal-core/built-in/commands/templates"
	logging "minimal/minimal-core/built-in/internal-logging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/stores/directory"
	"minimal/minimal-core/built-in/user-interfaces/tui"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"
)

func registerCommands(sourceGen *logging.SourceGenerator, commands *startup.Commands) {
	logger, _ := sourceGen.GetLogger("commandRegistry")
	
	err := commands.AddCommand("hello", helloWorld)

	if err != nil {
		logger.Fatal().Msg("Failed to register the hello command")
	}

	projectCreator := templates.NewProjectCreator(sourceGen)
	err = projectCreator.RegisterTemplateStore(directory.NewDirectoryStore(projectCreator.SourceGen))

	if err != nil {
		logger.Fatal().Msg("Failed to register the directory template store")
	}

	err = commands.AddCommand("new", projectCreator.NewProject)

	if err != nil {
		logger.Fatal().Msg("Failed to register the new command")
	}

	err = commands.AddCommand("message", func(_ []string) bool { return messagingTest(sourceGen) })

	if err != nil {
		logger.Fatal().Msg("Failed to register the message command")
	}

	err = commands.AddCommand("tui", tui.NewTUI().StartTUI)

	if err != nil {
		logger.Fatal().Msg("Failed to register the tui command")
	}
}

func helloWorld(_ []string) bool {
	fmt.Println("Hello, world!")
	return true
}

func messagingTest(sourceGen *logging.SourceGenerator) bool {
	logRenderer := logrendering.NewLogRenderer(sourceGen, os.Stdout, logrendering.Config{})
	messaging := usermessaging.NewMessenger(sourceGen)
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
