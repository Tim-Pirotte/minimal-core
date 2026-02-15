package templates

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
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

func isWindowsPrivilegeError(err error) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	return errors.Is(err, syscall.Errno(1314))
}

func TestLoadTemplateSymlink(t *testing.T) {
	osFs := afero.NewOsFs()

	workDir, err := afero.TempDir(osFs, "", "template-test")
	
	if err != nil {
		t.Fatal(err)
	}

	defer osFs.RemoveAll(workDir)

	templateRoot := filepath.Join(workDir, "templates")
	templateName := "go-cli"
	sourceDir := filepath.Join(templateRoot, templateName)

	if err := osFs.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}

	originalFile := filepath.Join(sourceDir, "config.yaml")

	if err := afero.WriteFile(osFs, originalFile, []byte("port: 8080"), 0o644); err != nil {
		t.Fatal(err)
	}

	linker := osFs.(afero.Linker)

	linkPath := filepath.Join(sourceDir, "config-link.yaml")

	if err := linker.SymlinkIfPossible("config.yaml", linkPath); err != nil {
        if isWindowsPrivilegeError(err) {
            t.Skip("skipping: symlink privilege not held")
        }

		t.Fatalf("failed to create symlink: %v", err)
	}

	targetPath := filepath.Join(workDir, "output")

	if err := loadTemplate(osFs, sourceDir, targetPath); err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}

	copiedLink := filepath.Join(targetPath, "config-link.yaml")

	exists, err := afero.Exists(osFs, copiedLink)

	if err != nil || !exists {
		t.Fatalf("expected symlink %s to exist", copiedLink)
	}

	lstater := osFs.(afero.Lstater)

	info, isLstat, err := lstater.LstatIfPossible(copiedLink)

	if err != nil {
		t.Fatalf("lstat failed: %v", err)
	}

	if !isLstat || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected %s to be a symlink", copiedLink)
	}

	content, err := afero.ReadFile(osFs, copiedLink)
	if err != nil {
		t.Fatalf("failed reading through symlink: %v", err)
	}

	if string(content) != "port: 8080" {
		t.Errorf("unexpected content: got %s, want %s", string(content), "port: 8080")
	}
}

func TestLoadTemplateSymlinkFallbackCopiesContents(t *testing.T) {
	osFs := afero.NewOsFs()
	memFs := afero.NewMemMapFs()

	fs := afero.NewCopyOnWriteFs(osFs, memFs)

	workDir, err := afero.TempDir(osFs, "", "template-test")

	if err != nil {
		t.Fatal(err)
	}

	defer osFs.RemoveAll(workDir)

	templateRoot := filepath.Join(workDir, "templates")
	templateName := "go-cli"
	sourceDir := filepath.Join(templateRoot, templateName)

	if err := osFs.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}

	originalFile := filepath.Join(sourceDir, "config.yaml")

	if err := afero.WriteFile(osFs, originalFile, []byte("port: 8080"), 0o644); err != nil {
		t.Fatal(err)
	}

	linker := osFs.(afero.Linker)
	linkPath := filepath.Join(sourceDir, "config-link.yaml")

	err = linker.SymlinkIfPossible("config.yaml", linkPath)
	if err != nil {
		t.Skipf("skipping: cannot create symlink in test environment: %v", err)
	}

	targetPath := filepath.Join(workDir, "output")

	if err := loadTemplate(fs, sourceDir, targetPath); err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}

	copiedPath := filepath.Join(targetPath, "config-link.yaml")

	exists, err := afero.Exists(fs, copiedPath)

	if err != nil || !exists {
		t.Fatalf("expected %s to exist", copiedPath)
	}

	info, err := fs.Stat(copiedPath)

	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("expected fallback to regular file, but got symlink")
	}

	content, err := afero.ReadFile(fs, copiedPath)

	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}

	if string(content) != "port: 8080" {
		t.Fatalf("unexpected content: got %s, want %s", string(content), "port: 8080")
	}
}

func TestLoadTemplateDirSymlinkPreserved(t *testing.T) {
	osFs := afero.NewOsFs()

	workDir, err := afero.TempDir(osFs, "", "template-test")

	if err != nil {
		t.Fatal(err)
	}

	defer osFs.RemoveAll(workDir)

	sourceDir := filepath.Join(workDir, "templates", "go-cli")
	
	if err := osFs.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}

	realDir := filepath.Join(sourceDir, "assets")
	
	if err := osFs.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err = afero.WriteFile(osFs, filepath.Join(realDir, "file.txt"), []byte("hello"), 0o644)
	
	if err != nil {
		t.Fatal(err)
	}

	linker := osFs.(afero.Linker)
	linkPath := filepath.Join(sourceDir, "assets-link")

	err = linker.SymlinkIfPossible("assets", linkPath)

	if err != nil {
		t.Skipf("skipping: cannot create directory symlink: %v", err)
	}

	targetPath := filepath.Join(workDir, "output")

	if err := loadTemplate(osFs, sourceDir, targetPath); err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}

	copiedLink := filepath.Join(targetPath, "assets-link")

	lstater := osFs.(afero.Lstater)
	info, isLstat, err := lstater.LstatIfPossible(copiedLink)

	if err != nil {
		t.Fatalf("lstat failed: %v", err)
	}

	if !isLstat || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected directory symlink to be preserved")
	}
}

func TestLoadTemplateDirSymlinkFallbackCopiesContents(t *testing.T) {
	osFs := afero.NewOsFs()
	memFs := afero.NewMemMapFs()

	fs := afero.NewCopyOnWriteFs(osFs, memFs)

	workDir, err := afero.TempDir(osFs, "", "template-test")
	
	if err != nil {
		t.Fatal(err)
	}

	defer osFs.RemoveAll(workDir)

	sourceDir := filepath.Join(workDir, "templates", "go-cli")
	if err := osFs.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}

	realDir := filepath.Join(sourceDir, "assets")

	if err := osFs.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err = afero.WriteFile(osFs, filepath.Join(realDir, "file.txt"), []byte("hello"), 0o644)

	if err != nil {
		t.Fatal(err)
	}

	linker := osFs.(afero.Linker)
	linkPath := filepath.Join(sourceDir, "assets-link")

	err = linker.SymlinkIfPossible("assets", linkPath)
	
	if err != nil {
		t.Skipf("skipping: cannot create directory symlink: %v", err)
	}

	targetPath := filepath.Join(workDir, "output")

	if err := loadTemplate(fs, sourceDir, targetPath); err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}

	copiedPath := filepath.Join(targetPath, "assets-link")

	info, err := fs.Stat(copiedPath)

	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if info.Mode() & os.ModeSymlink != 0 {
		t.Fatalf("expected fallback to real directory, but got symlink")
	}

	contentPath := filepath.Join(copiedPath, "file.txt")
	content, err := afero.ReadFile(fs, contentPath)

	if err != nil {
		t.Fatalf("expected copied directory contents, but file missing: %v", err)
	}

	if string(content) != "hello" {
		t.Fatalf("unexpected content: got %s, want hello", string(content))
	}
}
