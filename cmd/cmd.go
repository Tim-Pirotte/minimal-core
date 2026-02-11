package main

import "minimal/minimal-core/built-in/startup"

func main() {
	commands := startup.NewCommands()
	startup.Start(commands)
}
