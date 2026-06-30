package cli

import (
	"bufio"
	"fmt"
	"io"
	messaging "minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"strings"
)

const RingBufferSize = 4096

type CLI struct {
    messenger     *messaging.Messenger
    inputReader   *bufio.Reader
    outputWriter  *bufio.Writer
}

func NewCli(
    messenger *messaging.Messenger,
    reader io.Reader,
    writer io.Writer,
) *CLI {
    messenger.AddOutput(logrendering.NewLogRenderer(os.Stdout))

    return &CLI{messenger, bufio.NewReader(reader), bufio.NewWriter(writer)}
}

func (c *CLI) PromptBool(question string, defaultTrue bool) (answer bool, ok bool) {
    sb := strings.Builder{}
    sb.WriteString(question)
    sb.WriteByte('(')

    if defaultTrue {
        sb.WriteByte('Y')
    } else {
        sb.WriteByte('y')
    }

    sb.WriteByte('/')

    if !defaultTrue {
        sb.WriteByte('N')
    } else {
        sb.WriteByte('n')
    }

    sb.WriteByte(')')

    q := sb.String()

    for {
        c.outputWriter.WriteString(q)

        text, err := c.inputReader.ReadString('\n')

        if err != nil {
            c.messenger.Send(messaging.Message{
                Reference: "TODO",
                Message: "Could not read a string from the provided input",
                Severity: messaging.Error,
            })

            return false, false
        }

        input := strings.TrimSpace(strings.ToLower(text))

        switch input {
        case "":
            return defaultTrue, true
        case "y", "yes", "true", "1":
            return true, true
        case "n", "no", "false", "0":
            return false, true
        default:
            fmt.Println("Please enter yes or no (y/n).")
        }
    }
}

func (c *CLI) PromptString(question, suggestion string) (answer string, ok bool) {
    sb := strings.Builder{}
    sb.WriteString(question)

    if suggestion != "" {
        sb.WriteString("(default: ")
        sb.WriteString(suggestion)
        sb.WriteByte(')')
    }

    c.outputWriter.WriteString(sb.String())

    text, err := c.inputReader.ReadString('\n')

    if err != nil {
        c.messenger.Send(messaging.Message{
            Reference: "TODO",
            Message: "Could not read a string from the provided input",
            Severity: messaging.Error,
        })

        return "", false
    }

    input := strings.TrimSpace(text)

    if input == "" {
        return suggestion, true
    }

    return input, true
}

func (c *CLI) HandleCrash() {}
