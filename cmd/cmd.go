package main

import (
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
	logRenderer := logrendering.NewLogRenderer(sourceGen, os.Stdout, logrendering.Config{})
	messaging := usermessaging.NewMessenger(sourceGen)
	defer messaging.Close()
	messaging.AddOutput(logRenderer)
	
	commands := startup.NewCommands(sourceGen, messaging)

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
