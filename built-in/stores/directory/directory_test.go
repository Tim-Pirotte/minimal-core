package directory

import (
	"minimal/minimal-lang/built-in/messenger"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTemplateFromDirectory(t *testing.T) {
	templateRoot, err := os.MkdirTemp("", "templates")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := os.RemoveAll(templateRoot)

		if err != nil {
			t.Error("Removing temporary template files failed")
		}
	}()

	templateName := "go-cli"
	sourceDir := filepath.Join(templateRoot, templateName)

	files := map[string]string{
		filepath.Join(sourceDir, "main.go"):          "package main",
		filepath.Join(sourceDir, "scripts/build.sh"): "#!/bin/bash",
	}

	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	targetPath, err := os.MkdirTemp("", "my-new-app")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := os.RemoveAll(templateRoot)

		if err != nil {
			t.Error("Removing temporary template files failed")
		}
	}()

	directoryStore := NewDirectoryStore(messenger.New())

	directoryStore.LoadTemplate(templateName, "cli-app", targetPath, nil)

	testCases := []struct {
		name            string
		expectedPath    string
		expectedContent string
	}{
		{"Root file", filepath.Join(targetPath, "main.go"), "package main"},
		{"Nested file", filepath.Join(targetPath, "scripts/build.sh"), "#!/bin/bash"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.expectedPath)
			if err != nil || info.IsDir() {
				t.Errorf("expected path %s does not exist or is a directory", tc.expectedPath)
				return
			}

			content, err := os.ReadFile(tc.expectedPath)
			if err != nil {
				t.Fatalf("could not read file %s: %v", tc.expectedPath, err)
			}

			if string(content) != tc.expectedContent {
				t.Errorf("content mismatch for %s: got %s, want %s", tc.expectedPath, string(content), tc.expectedContent)
			}
		})
	}
}
