package templates

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestLoadTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

    templateRoot := "/usr/local/share/templates"
    templateName := "go-cli"
    sourceDir := filepath.Join(templateRoot, templateName)

    files := map[string]string{
        filepath.Join(sourceDir, "main.go"):          "package main",
        filepath.Join(sourceDir, "scripts/build.sh"): "#!/bin/bash",
    }

    for path, content := range files {
        err := fs.MkdirAll(filepath.Dir(path), 0755)
        if err != nil {
            t.Fatal(err)
        }
        err = afero.WriteFile(fs, path, []byte(content), 0644)
        if err != nil {
            t.Fatal(err)
        }
    }

    targetPath := "/home/user/my-new-app"

	sourcePath := filepath.Join(templateRoot, templateName)

    err := loadTemplate(fs, sourcePath, targetPath)
    if err != nil {
        t.Fatalf("loadTemplate failed: %v", err)
    }

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
            exists, _ := afero.Exists(fs, tc.expectedPath)
            if !exists {
                t.Errorf("expected path %s does not exist", tc.expectedPath)
                return
            }

            content, err := afero.ReadFile(fs, tc.expectedPath)
            if err != nil {
                t.Fatalf("could not read file %s: %v", tc.expectedPath, err)
            }

            if string(content) != tc.expectedContent {
                t.Errorf("content mismatch for %s: got %s, want %s", tc.expectedPath, string(content), tc.expectedContent)
            }
        })
    }
}

func TestLoadTemplateSymLink(t *testing.T) {
	
}
