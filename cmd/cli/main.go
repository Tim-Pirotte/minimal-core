package main

import (
	"bufio"
	messaging "minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/startup"
	"minimal/minimal-core/built-in/user-interfaces/cli"
	"minimal/minimal-core/setup"
	"os"
)

func main() {
    os.Exit(run())
}

func run() int {
    messenger := messaging.NewMessenger()
    defer messenger.Close()

    logRenderer := logrendering.NewLogRenderer(os.Stdout)
    messenger.AddOutput(logRenderer)

    cli := cli.NewCli(messenger, bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))

    commands := startup.NewCommands(messenger)
    setup.RegisterCommands(commands, messenger, cli)

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
