package cli

import (
	"bufio"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func setupTestCLI(input string) *CLI {
	return &CLI{
		loggingBuffer: *logging.NewRingBuffer(256),
		inputReader: bufio.NewReader(strings.NewReader(input)),
		outputWriter: bufio.NewWriter(io.Discard),
		logger:      logging.Logger{Logger: zerolog.Nop()},
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

func TestHandleCrash(t *testing.T) {
    t.Run("successful_export", func(t *testing.T) {
        tmpDir := t.TempDir()
        
        input := "y\n" + tmpDir + "\n"
        c := setupTestCLI(input)
        
        testLog := "ERROR: something broke"
        c.loggingBuffer.Write([]byte(testLog))

        c.HandleCrash()

        expectedFile := filepath.Join(tmpDir, "crash_report.log")
        content, err := os.ReadFile(expectedFile)

        if err != nil {
            t.Fatalf("could not read exported file: %v", err)
        }

        if string(content) != testLog {
            t.Errorf("expected %q, got %q", testLog, string(content))
        }
    })

    t.Run("user_cancels", func(t *testing.T) {
        input := "n\n"
        c := setupTestCLI(input)
        
        c.HandleCrash()

        if _, err := os.Stat("crash_report.log"); !os.IsNotExist(err) {
            t.Error("file should not have been created")
            os.Remove("crash_report.log")
        }
    })

    t.Run("invalid_path_error", func(t *testing.T) {
        input := "y\n/non/existent/path/that/should/fail\n"
        c := setupTestCLI(input)
        
        c.HandleCrash()
    })
}
