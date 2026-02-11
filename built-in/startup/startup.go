package startup

import (
	"minimal/minimal-core/built-in/config"
	"os"
	"path/filepath"
)

type Commands struct {
	commands map[string]func()
}

func NewCommands() Commands {
	return Commands{make(map[string]func())}
}

func (c *Commands) AddCommand(name string, function func()) {
	if _, ok := c.commands[name]; ok {
		// TODO log error
	}

	c.commands[name] = function
}

type StartupConfig struct {
	Command string `toml:"command"`
}

func Start(commands Commands) {
	if len(os.Args) != 2 {
		// TODO log error
	}

	configuration := os.Args[1]

	startupConfig := &StartupConfig{}

	file, err := os.ReadFile(filepath.Join(".", "commands", configuration))

	if err != nil {
		// TODO log error
	}

	config.LoadConfig(string(file), startupConfig)

	if startupFunc, ok := commands.commands[startupConfig.Command]; ok {
		startupFunc()
	} else {
		// TODO log error
	}
}
