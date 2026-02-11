package startup

import (
	"minimal/minimal-core/built-in/config"
	"os"
	"path/filepath"
)

type Commands struct {
	commands map[string]func()
	configlessCommands map[string]func()
}

func NewCommands() Commands {
	return Commands{make(map[string]func()), make(map[string]func())}
}

func (c *Commands) AddCommand(name string, function func()) {
	if _, ok := c.commands[name]; ok {
		// TODO log error
	}

	c.commands[name] = function
}

func (c *Commands) AddConfiglessCommand(name string, function func()) {
	if _, ok := c.commands[name]; ok {
		// TODO log error
	}

	c.configlessCommands[name] = function
}

type StartupConfig struct {
	Command string `toml:"command"`
}

func Start(commands Commands) {
	if len(os.Args) < 2 {
		// TODO log error
	}

	configOrCommand := os.Args[1]

	if startupFunc, ok := commands.configlessCommands[configOrCommand]; ok {
		startupFunc()
		return
	}

	startupConfig := &StartupConfig{}

	file, err := os.ReadFile(filepath.Join(".", "commands", configOrCommand))

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
