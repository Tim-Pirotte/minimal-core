package directory

import (
	"io/fs"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

const (
	defaultTemplateName = "default"
	templatesFolderName = "templates"
	defaultTargetPath = "."
)

type DirectoryStore struct {
	logger zerolog.Logger
}

func NewDirectoryStore(sourceGen *logging.SourceGenerator) DirectoryStore {
	logger, _ := sourceGen.GetLogger("directory")
	
	return DirectoryStore{logger}
}

func (d *DirectoryStore) HasTemplate(name string) bool {
	panic("unimplemented")
}

func (d *DirectoryStore) LoadTemplate(name, projectName, destinationPath string) (ok bool) {
	if name == "" {
		name = defaultTemplateName
	}

	if destinationPath == "" {
		destinationPath = defaultTargetPath
	}

	executablePath, err := os.Executable()

	if err != nil {
		// TODO log error
	}

	sourcePath := filepath.Join(filepath.Dir(executablePath), templatesFolderName, defaultTemplateName)

	err = os.CopyFS(destinationPath, os.DirFS(sourcePath))
    
	switch err.(type) {
	case nil:
	case *fs.PathError:
		d.checkPath(sourcePath, "source_path")
		d.checkPath(destinationPath, "target_path")

		return false
	default:
		d.logger.Error().
			Err(err).
			Str("source_path", sourcePath).
			Str("target_path", destinationPath).
			Msg("unforseen error loading template")
	
		return false
	}

	return true
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

func (d *DirectoryStore) StoreTemplate(name, location string) (ok bool) {
	panic("unimplemented")
}
