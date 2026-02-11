package startup

import (
	"minimal/minimal-core/built-in/config"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type Commands struct {
	commands map[string]func()
	configlessCommands map[string]func()
	logger zerolog.Logger
}

func NewCommands(logger zerolog.Logger) Commands {
	return Commands{make(map[string]func()), make(map[string]func()), logger}
}

func (c *Commands) AddCommand(name string, function func()) {
	if _, ok := c.commands[name]; ok {
		c.logger.Fatal().
			Str("command_name", name).
			Msg("duplicate command name")
	}

	c.commands[name] = function

	c.logger.Info().
		Str("command_name", name).
		Msg("command registered")
}

func (c *Commands) AddConfiglessCommand(name string, function func()) {
	if _, ok := c.configlessCommands[name]; ok {
		c.logger.Fatal().
			Str("command_name", name).
			Msg("duplicate configless command name")
	}

	c.configlessCommands[name] = function
	
	c.logger.Info().
		Str("command_name", name).
		Msg("configless command registered")
}

type StartupConfig struct {
	Command string `toml:"command"`
}

const minimumExpectedArgs = 2

func (c *Commands) Start() {
	if len(os.Args) < minimumExpectedArgs {
		c.logger.Fatal().
			Int("min_expected_args", minimumExpectedArgs).
			Int("actual_args", len(os.Args)).
			Msg("not enough arguments")
	}

	configOrCommand := os.Args[1]

	if startupFunc, ok := c.configlessCommands[configOrCommand]; ok {
		c.logger.Info().
			Str("command_name", configOrCommand).
			Msg("running configless command")

		startupFunc()
		return
	}

	startupConfig := &StartupConfig{}

	file, err := os.ReadFile(filepath.Join(".", "commands", configOrCommand))

	if err != nil {
		c.logger.Fatal().
			Str("config_name", configOrCommand).
			Msg("config not found on file system")
	}

	config.LoadConfig(string(file), startupConfig)

	if startupFunc, ok := c.commands[startupConfig.Command]; ok {
		c.logger.Info().
			Str("command_name", startupConfig.Command).
			Str("config_name", configOrCommand).
			Msg("running config")

		startupFunc()
	} else {
		c.logger.Fatal().
			Str("command_name", startupConfig.Command).
			Str("config_name", configOrCommand).
			Msg("command does not exist")
	}
}
