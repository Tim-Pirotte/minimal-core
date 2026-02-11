package main

import (
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/startup"
	"os"
)

func main() {
	sourceGen := logging.Init(os.Stdout)
	cmdLogger, _ := sourceGen.GetLogger("commands")
	commands := startup.NewCommands(cmdLogger)

	registerCommands(commands)

	commands.Start()
}
