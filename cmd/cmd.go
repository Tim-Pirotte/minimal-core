package main

import (
	"minimal/minimal-core/built-in/config-loaders/shell"
	"minimal/minimal-core/built-in/config-loaders/toml"
	logging "minimal/minimal-core/built-in/internal-logging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/startup"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	os.Exit(run())
}

func run() int {
	sourceGen := logging.Init(zerolog.ConsoleWriter{Out: os.Stdout})

	configLoader := shell.NewShell(toml.NewConfigLoader(sourceGen))

	logRenderer := logrendering.NewLogRenderer(sourceGen, os.Stdout, configLoader)
	messaging := usermessaging.NewMessenger(sourceGen)
	defer messaging.Close()
	messaging.AddOutput(logRenderer)
	
	commands := startup.NewCommands(sourceGen, messaging, configLoader)

	registerCommands(sourceGen, commands)

	entrypoint, args := commands.GetEntrypoint(os.Args)

	if entrypoint == nil {
		return 1
	}

	ok := entrypoint(args)

	if !ok {
		return 1
	}

	return 0
}
