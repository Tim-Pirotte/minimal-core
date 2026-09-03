package main

import (
	"bufio"
	messaging "minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/outputs/log-renderer"
	"minimal/minimal-lang/built-in/startup"
	"minimal/minimal-lang/built-in/user-interfaces/cli"
	"minimal/minimal-lang/setup"
	"os"
)

func main() {
    os.Exit(run())
}

func run() int {
    messenger := messaging.New()
    defer messenger.Close()

    logRenderer := logrenderer.New(os.Stdout)
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
