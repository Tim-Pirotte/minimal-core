package main

import (
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	sourceGen := logging.Init(zerolog.ConsoleWriter{Out: os.Stdout})
	commands := startup.NewCommands(sourceGen)

	registerCommands(commands)

	entrypoint := commands.GetEntrypoint(os.Args)

	if entrypoint == nil {
		os.Exit(1)
	}

	entrypoint()
}
