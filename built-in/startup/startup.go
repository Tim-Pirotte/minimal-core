package startup

import (
	"errors"
	"minimal/minimal-core/built-in/config"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

const minimumExpectedArgs = 2

var DuplicateCommand = errors.New("command with this name already exists")
var commandsConfigPath = filepath.Join(".", "commands")

type Commands struct {
	commands map[string]func()
	logger zerolog.Logger
}

func NewCommands(logger zerolog.Logger) *Commands {
	return &Commands{make(map[string]func()), logger}
}

func (c *Commands) AddCommand(name string, function func()) error {
	if _, ok := c.commands[name]; ok {
		c.logDuplicateCommand(name)
		return DuplicateCommand
	}

	c.commands[name] = function
	c.logCommandRegistered(name)

	return nil
}

type StartupConfig struct {
	Command string `toml:"command"`
}

func (c *Commands) GetEntrypoint() func() {
	if len(os.Args) < minimumExpectedArgs {
		c.logNotEnoughArgs()
		panic("unreachable")
	}

	configOrCommand := os.Args[1]

	if startupFunc, ok := c.commands[configOrCommand]; ok {
		c.logRunningCommand(configOrCommand, false)
		return startupFunc
	} else {
		return c.loadFromConfig(configOrCommand)
	}
}

func (c *Commands) loadFromConfig(configName string) func() {
	startupConfig := &StartupConfig{}

	file, err := os.ReadFile(filepath.Join(commandsConfigPath, configName))

	if err != nil {
		c.logConfigNotFound(configName)
		panic("unreachable")
	}

	config.LoadConfig(string(file), startupConfig)

	if startupFunc, ok := c.commands[startupConfig.Command]; ok {
		c.logRunningCommand(startupConfig.Command, true)
		return startupFunc
	}
	
	c.logCommandNotExists(startupConfig.Command, configName)
	panic("unreachable")
}

func (c *Commands) logDuplicateCommand(name string) {
	c.logger.Error().
		Str("command_name", name).
		Err(DuplicateCommand).
		Msg("")
}

func (c *Commands) logCommandRegistered(name string) {
	c.logger.Debug().
		Str("command_name", name).
		Msg("command registered")
}

func (c *Commands) logNotEnoughArgs() {
	c.logger.Fatal().
		Int("min_expected_args", minimumExpectedArgs).
		Int("actual_args", len(os.Args)).
		Msg("not enough arguments")
}

func (c *Commands) logRunningCommand(commandName string, fromConfig bool) {
	c.logger.Info().
		Str("command_name", commandName).
		Bool("from_config", fromConfig).
		Msg("running command")
}

func (c *Commands) logConfigNotFound(configName string) {
	c.logger.Fatal().
		Str("config_name", configName).
		Msg("config not found on file system")
}

func (c *Commands) logCommandNotExists(commandName, configName string) {
	c.logger.Fatal().
		Str("command_name", commandName).
		Str("config_name", configName).
		Msg("command does not exist")
}
