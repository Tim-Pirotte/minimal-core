package templating

import (
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	defaultTemplateName = "default"
	templatesFolderName = "templates"

	destinationFlagName = "destination"
)

func NewProject() {
	var targetPath string
	flag.StringVar(&targetPath, destinationFlagName, "", "")
	flag.StringVar(&targetPath, string(destinationFlagName[0]), "", "")

	var symbolicLink bool
	flag.BoolVar(&symbolicLink, "ln", false, "")

	flag.Parse()

	if targetPath == "" {
		targetPath = "."
	}

	switch len(flag.Args()) {
	case 0:
		loadTemplate(targetPath, defaultTemplateName)
	case 1:
		loadTemplate(targetPath, flag.Arg(0))
	default:
		// TODO log error
	}
}

func loadTemplate(targetPath, name string) {
	executablePath, err := os.Executable()

	if err != nil {
		// TODO log error
	}

	templatePath := filepath.Join(filepath.Dir(executablePath), templatesFolderName, name)

	err = filepath.WalkDir(templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// TODO check if the error is of importance
		relativePath, _ := filepath.Rel(templatePath, path)
        destinationPath := filepath.Join(targetPath, relativePath)

		if d.IsDir() {
			return os.MkdirAll(destinationPath, 0755)
		} else {
			sourceFile, err := os.Open(path)

			if err != nil {
				return err
			}

			defer sourceFile.Close()

			destinationFile, err := os.Create(destinationPath)

			if err != nil {
				return err
			}

			defer destinationFile.Close()

			_, err = io.Copy(destinationFile, sourceFile)

			return err
		}
	})
}
