package cli

import (
	"bufio"
	"fmt"
	configloader "minimal/minimal-core/built-in/config-loader"
	logging "minimal/minimal-core/built-in/internal-logging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

const RingBufferSize = 4096

type CLI struct {
	loggingBuffer logging.RingBuffer
	inputReader   *bufio.Reader
	outputWriter  *bufio.Writer
	logger        zerolog.Logger
}

func NewCli(
	sourceGen *logging.SourceGenerator, 
	messenger *usermessaging.Messenger, 
	reader *bufio.Reader, 
	writer *bufio.Writer,
	configLoader configloader.ConfigLoader,
) *CLI {
	messenger.AddOutput(logrendering.NewLogRenderer(sourceGen, os.Stdout, configLoader))
	logger, _ := sourceGen.GetLogger("cli")

	return &CLI{*logging.NewRingBuffer(RingBufferSize), reader, writer, logger}
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
		c.logger.Error().Err(err).Msg("string read")
		return "", false
	}

	input := strings.TrimSpace(text)

	if input == "" {
		return suggestion, true
	}

	return input, true
}

func (c *CLI) HandleCrash() {
	exportLog, ok := c.PromptBool("Would you like to export the logs for a bug report?", false)

	if !(ok && exportLog) { return }

	path, ok := c.PromptString("Where would you like to have this exported?", ".")

	if !ok { return }

	fullPath := filepath.Join(path, "crash_report.log")

	f, err := os.Create(fullPath)

    if err != nil {
        c.logger.Error().Err(err).Str("path", fullPath).Msg("crash report create file failed")
        return
    }

	defer f.Close()

	_, err = c.loggingBuffer.WriteTo(f)

	if err != nil {
		c.logger.Error().Err(err).Str("path", fullPath).Msg("copying ring buffer failed")
		return
	}

	c.logger.Debug().Msg("exported internal logs")
}
