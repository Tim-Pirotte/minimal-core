package main

import "minimal/minimal-core/built-in/startup"

func main() {
	commands := startup.NewCommands()
	registerCommands(commands)
	startup.Start(commands)
}
