package main

import (
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	sourceGen := logging.Init(zerolog.ConsoleWriter{Out: os.Stdout})
	cmdLogger, _ := sourceGen.GetLogger("commands")
	commands := startup.NewCommands(cmdLogger)

	registerCommands(commands)

	entrypoint := commands.GetEntrypoint()

	if entrypoint == nil {
		os.Exit(1)
	}
	
	entrypoint()
}
