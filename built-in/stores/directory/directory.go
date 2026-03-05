package directory

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

const (
	defaultStoreName = "directory"
	defaultTemplateName = "default"
	templatesFolderName = "templates"
	defaultTargetPath = "."
)

type DirectoryStore struct {
	StoreName string
	logger zerolog.Logger
}

func NewDirectoryStore(sourceGen *logging.SourceGenerator) *DirectoryStore {
	logger, _ := sourceGen.GetLogger(defaultStoreName)
	
	return &DirectoryStore{defaultStoreName, logger}
}

func (d *DirectoryStore) Name() string {
	return d.StoreName
}

func (d *DirectoryStore) HasTemplate(name string) bool {
	if name == "" {
		name = defaultTemplateName
	}

	sourcePath, ok := d.getSourcePath(name)

	if !ok {
		return false
	}

	_, err := os.Stat(sourcePath)

	return err == nil
}

func (d *DirectoryStore) LoadTemplate(name, projectName, destinationPath string, args []string) (ok bool) {
	fmt.Println(args)

	if name == "" {
		name = defaultTemplateName
	}

	if destinationPath == "" {
		destinationPath = defaultTargetPath
	}

	sourcePath, ok := d.getSourcePath(name)

	if !ok {
		return false
	}

	targetPath := filepath.Join(destinationPath, projectName)

	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return false
	}

	err := os.CopyFS(targetPath, os.DirFS(sourcePath))
    
	switch err.(type) {
	case nil:
	case *fs.PathError:
		d.checkPath(sourcePath, "source_path")
		d.checkPath(targetPath, "target_path")

		return false
	default:
		d.logger.Error().
			Err(err).
			Str("source_path", sourcePath).
			Str("target_path", targetPath).
			Msg("unforseen error loading template")
	
		return false
	}

	return true
}

func (d *DirectoryStore) getSourcePath(name string) (string, bool) {
	executablePath, err := os.Executable()

	if err != nil {
		d.logger.Error().Str("project_name", name).Msg("cannot retrieve the executable path")
		return "", false
	}

	return filepath.Join(filepath.Dir(executablePath), templatesFolderName, name), true
}

func (d *DirectoryStore) checkPath(path, reference string) {
	_, err := os.Stat(path)
	
	if err != nil {
		d.logger.Error().
			Str("path", path).
			Str("reference", reference).
			Msg("path does not exist")
	}
}

func (d *DirectoryStore) StoreTemplate(name, location, args []string) (ok bool) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // TODO log error
	
	var symbolicLink bool
	fs.BoolVar(&symbolicLink, "ln", false, "")

	if err := fs.Parse(args); err != nil {
        return false
    }

	panic("unimplemented")
}
