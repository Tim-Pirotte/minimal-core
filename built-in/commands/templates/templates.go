package templates

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	defaultTemplateName = "default"
	templatesFolderName = "templates"

	destinationFlagName = "destination"
	defaultTargetPath = "."
)

func NewProject() {
	var targetPath string
	flag.StringVar(&targetPath, destinationFlagName, defaultTargetPath, "")
	flag.StringVar(&targetPath, string(destinationFlagName[0]), defaultTargetPath, "")

	flag.Parse()

	executablePath, err := os.Executable()

	if err != nil {
		// TODO log error
		err = nil
	}

	templatePath := filepath.Join(filepath.Dir(executablePath), templatesFolderName)

	fileSystem := afero.NewOsFs()

	switch flag.NArg() {
	case 0:
		sourcePath := filepath.Join(templatePath, defaultTemplateName)
		err = loadTemplate(fileSystem, sourcePath, targetPath)
	case 1:
		sourcePath := filepath.Join(templatePath, flag.Arg(0))
		err = loadTemplate(fileSystem, sourcePath, targetPath)
	default:
		// TODO log error
	}

	if err != nil {

	}
}

func loadTemplate(fs afero.Fs, sourcePath, targetPath string) error {
    return afero.Walk(fs, sourcePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        relPath, _ := filepath.Rel(sourcePath, path)
        destPath := filepath.Join(targetPath, relPath)

        if info.IsDir() {
            return fs.MkdirAll(destPath, 0755)
        } else {
			srcFile, err := fs.Open(path)

			if err != nil {
				return err
			}

			defer srcFile.Close()

			dstFile, err := fs.Create(destPath)

			if err != nil {
				return err
			}

			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			return err
		}
    })
}

func CreateTemplate() {
	var symbolicLink bool
	flag.BoolVar(&symbolicLink, "ln", false, "")

	flag.Parse()
}

func saveTemplate() {

}
