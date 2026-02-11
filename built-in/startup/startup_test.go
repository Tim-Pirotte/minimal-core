package startup

import (
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
)

func setupTestCommands(mockFiles fstest.MapFS) *Commands {
	return &Commands{
		commands: make(map[string]func()),
		logger:   zerolog.Nop(),
		fs:       mockFiles,
	}
}

func TestAddCommand(t *testing.T) {
	c := setupTestCommands(nil)
	err := c.AddCommand("build", func() {})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, ok := c.commands["build"]; !ok {
		t.Error("command was not actually registered in the map")
	}
}

func TestAddDuplicateCommand(t *testing.T) {
	c := setupTestCommands(nil)
	_ = c.AddCommand("build", func() {})
	err := c.AddCommand("build", func() {})

	if err != DuplicateCommand {
		t.Errorf("expected DuplicateCommand error, got %v", err)
	}
}

func TestGetEntrypointFromCommand(t *testing.T) {
	c := setupTestCommands(nil)
	called := false
	c.AddCommand("build", func() { called = true })

	fn := c.GetEntrypoint([]string{"app", "build"})

	if fn == nil {
		t.Fatal("expected function, got nil")
	}

	fn()

	if !called {
		t.Error("the returned function was not the registered one")
	}
}

func TestInvalidEntrypointFromCommand(t *testing.T) {
	c := setupTestCommands(fstest.MapFS{})
	
	fn := c.GetEntrypoint([]string{"app", "not-here"})

	if fn != nil {
		t.Error("expected nil for non-existent command/config")
	}
}

func TestGetEntrypointFromConfig(t *testing.T) {
	mockFS := fstest.MapFS{
		"commands/my-setup.toml": {
			Data: []byte(`command = "compile"`),
		},
	}

	c := setupTestCommands(mockFS)
	called := false
	c.AddCommand("compile", func() { called = true })

	fn := c.GetEntrypoint([]string{"app", "my-setup.toml"})

	if fn == nil {
		t.Fatal("expected function from config, got nil")
	}

	fn()

	if !called {
		t.Error("the function from config was not executed")
	}
}

func TestNonexistentConfig(t *testing.T) {
	c := setupTestCommands(fstest.MapFS{})
	
	fn := c.GetEntrypoint([]string{"app", "ghost.toml"})

	if fn != nil {
		t.Error("expected nil when config file is missing")
	}
}

func TestNonexistentCommandInConfig(t *testing.T) {
	mockFS := fstest.MapFS{
		"commands/bad-config.toml": {
			Data: []byte(`command = "missing-cmd"`),
		},
	}
	c := setupTestCommands(mockFS)

	fn := c.GetEntrypoint([]string{"app", "bad-config.toml"})

	if fn != nil {
		t.Error("expected nil because the command inside the config is not registered")
	}
}
