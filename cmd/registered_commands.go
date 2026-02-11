package main

import (
	"fmt"
	"minimal/minimal-core/built-in/startup"
)

func registerCommands(commands startup.Commands) {
	commands.AddConfiglessCommand("hello", helloWorld)
}

func helloWorld() {
	fmt.Println("Hello, world!")
}
