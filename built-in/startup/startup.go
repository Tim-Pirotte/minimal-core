package startup

import (
	"errors"
	"minimal/minimal-core/built-in/config"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

var DuplicateCommand = errors.New("command with this name already exists")

type Commands struct {
	commands map[string]func()
	configlessCommands map[string]func()
	logger zerolog.Logger
}

func NewCommands(logger zerolog.Logger) Commands {
	return Commands{make(map[string]func()), make(map[string]func()), logger}
}

func (c *Commands) AddCommand(name string, function func()) error {
	if _, ok := c.commands[name]; ok {
		c.log_duplicate_command(name)
		return DuplicateCommand
	}

	if _, ok := c.configlessCommands[name]; ok {
		c.log_duplicate_command(name)
		return DuplicateCommand
	}

	c.commands[name] = function

	c.log_command_registered(name, false)

	return nil
}

func (c *Commands) AddConfiglessCommand(name string, function func()) error {
	if _, ok := c.configlessCommands[name]; ok {
		c.log_duplicate_configless_command(name)
		return DuplicateCommand
	}

	if _, ok := c.commands[name]; ok {
		c.log_duplicate_configless_command(name)
		return DuplicateCommand
	}

	c.configlessCommands[name] = function
	
	c.log_command_registered(name, true)

	return nil
}

type StartupConfig struct {
	Command string `toml:"command"`
}

const minimumExpectedArgs = 2

func (c *Commands) Start() {
	if len(os.Args) < minimumExpectedArgs {
		c.log_not_enough_args()
	}

	configOrCommand := os.Args[1]

	if startupFunc, ok := c.configlessCommands[configOrCommand]; ok {
		c.log_run_configless_command(configOrCommand)
		startupFunc()

		return
	}

	startupConfig := &StartupConfig{}

	file, err := os.ReadFile(filepath.Join(".", "commands", configOrCommand))

	if err != nil {
		c.log_config_not_found(configOrCommand)
	}

	config.LoadConfig(string(file), startupConfig)

	if startupFunc, ok := c.commands[startupConfig.Command]; ok {
		c.log_running_config(startupConfig.Command, configOrCommand)

		startupFunc()
	} else {
		c.log_command_not_exists(startupConfig.Command, configOrCommand)
	}
}

func (c *Commands) log_duplicate_command(name string) {
	c.logger.Fatal().
		Str("command_name", name).
		Msg("duplicate command name")
}

func (c *Commands) log_command_registered(name string, withConfig bool) {
	c.logger.Debug().
		Str("command_name", name).
		Bool("with_config", withConfig).
		Msg("command registered")
}

func (c *Commands) log_duplicate_configless_command(name string) {
	c.logger.Fatal().
		Str("command_name", name).
		Msg("duplicate configless command name")
}

func (c *Commands) log_not_enough_args() {
	c.logger.Fatal().
		Int("min_expected_args", minimumExpectedArgs).
		Int("actual_args", len(os.Args)).
		Msg("not enough arguments")
}

func (c *Commands) log_run_configless_command(commandName string) {
	c.logger.Info().
		Str("command_name", commandName).
		Msg("running configless command")
}

func (c *Commands) log_config_not_found(configName string) {
	c.logger.Fatal().
		Str("config_name", configName).
		Msg("config not found on file system")
}

func (c *Commands) log_running_config(commandName, configName string) {
	c.logger.Info().
		Str("command_name", commandName).
		Str("config_name", configName).
		Msg("running config")
}

func (c *Commands) log_command_not_exists(commandName, configName string) {
	c.logger.Fatal().
		Str("command_name", commandName).
		Str("config_name", configName).
		Msg("command does not exist")
}
