package cli

import (
	"bufio"
	"fmt"
	logrendering "minimal/minimal-core/built-in/extensions/outputs/log-renderer"
	logging "minimal/minimal-core/built-in/internal-logging"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const RingBufferSize = 4096

type CLI struct {
	loggingBuffer logging.RingBuffer
	inputReader   *bufio.Reader
	logger        zerolog.Logger
}

func NewCli(sourceGen *logging.SourceGenerator, messenger *usermessaging.Messenger, reader *bufio.Reader) *CLI {
	// TODO load config properly
	messenger.AddOutput(logrendering.NewLogRenderer(sourceGen, os.Stdout, logrendering.Config{}))
	logger, _ := sourceGen.GetLogger("cli")

	return &CLI{*logging.NewRingBuffer(RingBufferSize), reader, logger}
}

func (c *CLI) PromptBool(question string, defaultTrue bool) (bool, ok bool) {
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
		fmt.Print(q)
		
		text, err := c.inputReader.ReadString('\n')
		
		if err != nil {
			c.logger.Error().Err(err).Msg("bool read")
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

func (c *CLI) PromptString(question, suggestion string) string {
	sb := strings.Builder{}
	sb.WriteString(question)

	if suggestion != "" {
		sb.WriteString("(default: ")
		sb.WriteString(suggestion)
		sb.WriteByte(')')
	}

	fmt.Print(sb.String())

	text, err := c.inputReader.ReadString('\n')

	if err != nil {
		c.logger.Error().Err(err).Msg("string read")
		// TODO do something when this error occurs
	}

	input := strings.TrimSpace(text)

	if input == "" {
		return suggestion
	}

	return input
}
