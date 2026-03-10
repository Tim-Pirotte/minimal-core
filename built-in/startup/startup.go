package startup

import (
	"errors"
	"io/fs"
	"minimal/minimal-core/built-in/config"
	logging "minimal/minimal-core/built-in/internal-logging"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

const (
	minimumExpectedArgs = 2
	baseConfigName = "base.toml"
)

var ErrDuplicateCommand = errors.New("command with this name already exists")
var commandsConfigPath = path.Join(".", "commands")

type Commands struct {
	commands map[string]func(args []string) (ok bool)
	logger zerolog.Logger
	messenger *usermessaging.Messenger
	fs fs.FS
}

func NewCommands(sourceGen *logging.SourceGenerator, messenger *usermessaging.Messenger) *Commands {
	logger, _ := sourceGen.GetLogger("startup")

	return &Commands{
		make(map[string]func(args []string) (ok bool)), 
		logger,
		messenger,
		os.DirFS("."),
	}
}

func (c *Commands) AddCommand(name string, function func(args []string) (ok bool)) error {
	if _, ok := c.commands[name]; ok {
		c.logDuplicateCommand(name)
		return ErrDuplicateCommand
	}

	c.commands[name] = function
	c.logCommandRegistered(name)

	return nil
}

type StartupConfig struct {
	Command string `toml:"command"`
}

// Returns the program entrypoint based on the first argument
// or nil if something went wrong
func (c *Commands) GetEntrypoint(args []string) (fn func(args []string) (ok bool), arguments []string) {
	if len(args) < minimumExpectedArgs {
		c.logNotEnoughArgs(len(args))
		return nil, nil
	}

	configOrCommand := args[1]

	if startupFunc, ok := c.commands[configOrCommand]; ok {
		c.logRunningCommand(configOrCommand, false)
		return startupFunc, args[2:]
	} else {
		return c.loadFromConfig(configOrCommand), args[2:]
	}
}

func (c *Commands) loadFromConfig(configName string) func(args []string) (ok bool) {
	startupConfig := &StartupConfig{}

	file, err := fs.ReadFile(c.fs, path.Join(commandsConfigPath, configName, baseConfigName))

	if err != nil {
		c.logConfigNotFound(configName)
		return nil
	}

	err = config.LoadConfig(string(file), startupConfig)

	if err != nil {
		c.logger.Error().Err(err).Msg("error loading config")

		if parserErr, ok := err.(toml.ParseError); ok {
			t := c.messenger.CreateLogTransaction()

			c.messenger.LogMessage(
				t, 
				usermessaging.Message{
					Severity: usermessaging.Critical, 
					Category: "ConfigError", 
					Message: parserErr.Error(),
				},
			)

			c.messenger.CommitLogTransaction(t)
		}

		return nil
	}

	if startupFunc, ok := c.commands[startupConfig.Command]; ok {
		c.logRunningCommand(startupConfig.Command, true)
		return startupFunc
	}
	
	c.logCommandNotExists(startupConfig.Command, configName)

	return nil
}

func (c *Commands) logDuplicateCommand(name string) {
	c.logger.Error().
		Str("command_name", name).
		Err(ErrDuplicateCommand).
		Msg("")
}

func (c *Commands) logCommandRegistered(name string) {
	c.logger.Debug().
		Str("command_name", name).
		Msg("command registered")
}

func (c *Commands) logNotEnoughArgs(argsLength int) {
	c.logger.Error().
		Int("min_expected_args", minimumExpectedArgs).
		Int("actual_args", argsLength).
		Msg("not enough arguments")
}

func (c *Commands) logRunningCommand(commandName string, fromConfig bool) {
	c.logger.Info().
		Str("command_name", commandName).
		Bool("from_config", fromConfig).
		Msg("running command")
}

func (c *Commands) logConfigNotFound(configName string) {
	c.logger.Error().
		Str("config_name", configName).
		Msg("config not found on file system")
}

func (c *Commands) logCommandNotExists(commandName, configName string) {
	c.logger.Error().
		Str("command_name", commandName).
		Str("config_name", configName).
		Msg("command does not exist")
}
