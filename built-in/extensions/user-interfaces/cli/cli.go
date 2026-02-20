package cli

import (
	logging "minimal/minimal-core/built-in/internal-logging"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

type CLI struct {
	loggingBuffer logging.RingBuffer
}

func NewCli(usermessaging.Messenger) *CLI {
	return &CLI{}
}
