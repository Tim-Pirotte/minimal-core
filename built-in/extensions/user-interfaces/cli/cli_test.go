package cli

import (
	"bufio"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func setupTestCLI(input string) *CLI {
	return &CLI{
		inputReader: bufio.NewReader(strings.NewReader(input)),
		logger:      zerolog.Nop(),
	}
}

func TestPromptBool(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultTrue  bool
		wantAnswer   bool
		wantOk       bool
	}{
		{"Empty input returns default (true)", "\n", true, true, true},
		{"Empty input returns default (false)", "\n", false, false, true},
		{"Yes variations (y)", "y\n", false, true, true},
		{"Yes variations (yes)", "YES\n", false, true, true},
		{"Yes variations (1)", "1\n", false, true, true},
		{"No variations (n)", "n\n", true, false, true},
		{"No variations (no)", "no\n", true, false, true},
		{"No variations (0)", "0\n", true, false, true},
		{"Retry on invalid then yes", "maybe\nnot-sure\ny\n", false, true, true},
		{"Handle EOF/Error", "", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestCLI(tt.input)
			ans, ok := c.PromptBool("Continue?", tt.defaultTrue)
			if ans != tt.wantAnswer || ok != tt.wantOk {
				t.Errorf("PromptBool() = (%v, %v), want (%v, %v)", ans, ok, tt.wantAnswer, tt.wantOk)
			}
		})
	}
}

func TestPromptString(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		suggestion string
		wantAnswer string
		wantOk     bool
	}{
		{"Normal input", "Hello World\n", "default", "Hello World", true},
		{"Empty input returns suggestion", "\n", "default_val", "default_val", true},
		{"Trims whitespace", "  spaced input  \n", "", "spaced input", true},
		{"Handle EOF/Error", "", "suggestion", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestCLI(tt.input)
			ans, ok := c.PromptString("Enter text:", tt.suggestion)
			if ans != tt.wantAnswer || ok != tt.wantOk {
				t.Errorf("PromptString() = (%v, %v), want (%q, %v)", ans, ok, tt.wantAnswer, tt.wantOk)
			}
		})
	}
}
