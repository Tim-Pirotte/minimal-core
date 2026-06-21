package main

import (
	"bufio"
	"minimal/minimal-core/built-in/config-loaders/shell"
	"minimal/minimal-core/built-in/config-loaders/toml"
	logging "minimal/minimal-core/built-in/internal-logging"
	messaging "minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/user-interfaces/cli"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	os.Exit(run())
}

func run() int {
	sourceGen := logging.Init(zerolog.ConsoleWriter{Out: os.Stdout})

	configLoader := shell.NewShellConfigLoader(toml.NewConfigLoader(sourceGen))

	messenger := messaging.NewMessenger(sourceGen)
	defer messenger.Close()

	logRenderer := logrendering.NewLogRenderer(sourceGen, os.Stdout, configLoader)
	messenger.AddOutput(logRenderer)

	cli := cli.NewCli(sourceGen, messenger, bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout), configLoader)

	commands := startup.NewCommands(sourceGen, messenger, configLoader)
	registerCommands(commands, sourceGen, messenger, cli)

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
