package templates

import (
	"flag"
	"os"
	"path/filepath"
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

	switch flag.NArg() {
	case 0:
		sourcePath := filepath.Join(templatePath, defaultTemplateName)
		err = loadTemplate(sourcePath, targetPath)
	case 1:
		sourcePath := filepath.Join(templatePath, flag.Arg(0))
		err = loadTemplate(sourcePath, targetPath)
	default:
		// TODO log error
	}

	if err != nil {

	}
}

func loadTemplate(sourcePath, targetPath string) error {
    return os.CopyFS(targetPath, os.DirFS(sourcePath))
}

func CreateTemplate() {
	var symbolicLink bool
	flag.BoolVar(&symbolicLink, "ln", false, "")

	flag.Parse()
}

func saveTemplate() {

}
