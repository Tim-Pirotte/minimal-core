package startup

import (
	"errors"
	"fmt"
	"io/fs"
	logging "minimal/minimal-core/built-in/internal-logging"
	messaging "minimal/minimal-core/built-in/messaging"
	"os"
	"path"
)

const minimumExpectedArgs = 2

var ErrDuplicateCommand = errors.New("command with this name already exists")
var commandsConfigPath = path.Join(".", "commands")

type Commands struct {
	commands map[string]func(args []string) (ok bool)
	logger logging.Logger
	messenger *messaging.Messenger
	fs fs.FS
}

func NewCommands(
	sourceGen *logging.SourceGenerator,
	messenger *messaging.Messenger,
) *Commands {
	logger, _ := sourceGen.GetLogger("Startup")

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

// Returns the program entrypoint based on the first argument or nil if something went wrong.
// The config gets replaced by no-config when a command is executed directly.
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
		// TODO what here?
		// return c.loadFromConfig(configOrCommand), args[2:]
	}

	panic("unimplemented")
}

// func (c *Commands) loadFromConfig(configName string) func(args []string) (ok bool) {
// 	c.ConfigLoader.SetLocalConfigSource(configName)

// 	commandAny, ok := c.ConfigLoader.Get("base", "execute", "command")

// 	if !ok {
// 		return nil
// 	}

// 	command, ok := commandAny.(string)

// 	if !ok {
// 		return nil
// 	}

// 	if startupFunc, ok := c.commands[command]; ok {
// 		c.logRunningCommand(command, true)
// 		return startupFunc
// 	}

// 	c.logCommandNotExists(command, configName)

// 	return nil
// }

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

	messageParts := []messaging.MessagePart{
		&messaging.Message{
			Severity: messaging.Critical,
			Message: fmt.Sprintf(
				"Expected at least %d arguments but got %d arguments",
				minimumExpectedArgs,
				argsLength,
			),
		},
		&messaging.Hint{Text: "mnm <commandName>", MoreInfoReference: ""},
	}

	c.messenger.Send(messageParts)
}

func (c *Commands) logRunningCommand(commandName string, fromConfig bool) {
	c.logger.Info().
		Str("command_name", commandName).
		Bool("from_config", fromConfig).
		Msg("running command")
}

func (c *Commands) logCommandNotExists(commandName, configName string) {
	c.logger.Error().
		Str("command_name", commandName).
		Str("config_name", configName).
		Msg("command does not exist")
}
