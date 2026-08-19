package startup

import (
	"errors"
	"io/fs"
	messaging "minimal/minimal-core/built-in/messenger"
	"os"
	"path"
	"strconv"
)

const minimumExpectedArgs = 2

var ErrDuplicateCommand = errors.New("command with this name already exists")
var commandsConfigPath = path.Join(".", "commands")

type Commands struct {
    commands map[string]func(args []string) (ok bool)
    messenger *messaging.Messenger
    fs fs.FS
}

func NewCommands(messenger *messaging.Messenger) *Commands {
    return &Commands{
        make(map[string]func(args []string) (ok bool)),
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
    c.messenger.Send(messaging.Message{
        Message: "The command '" + name + "' has already been declared",
        Severity: messaging.Error,
    })
}

func (c *Commands) logNotEnoughArgs(argsLength int) {
    c.messenger.Send(messaging.Message{
        Message: "Expected " +
                 strconv.Itoa(minimumExpectedArgs) +
                 " arguments but got " +
                 strconv.Itoa(argsLength),
        Severity: messaging.Error,
    })
}

func (c *Commands) logRunningCommand(commandName string, fromConfig bool) {
    c.messenger.Send(messaging.Message{
        Message: "Running '" + commandName + "'",
        Severity: messaging.Debug,
    })
}

func (c *Commands) logCommandNotExists(commandName, configName string) {
    c.messenger.Send(messaging.Message{
        Message: "The command '" + commandName + "' does not exist",
        Severity: messaging.Error,
    })
}
