package directory

import (
	"flag"
	"io"
	"io/fs"
	"minimal/minimal-core/built-in/messaging"
	"os"
	"path/filepath"
)

const (
    defaultTemplateName = "default"
    templatesFolderName = "templates"
    defaultTargetPath = "."
)

type DirectoryStore struct {
    messenger *messaging.Messenger
}

func NewDirectoryStore(messenger *messaging.Messenger) *DirectoryStore {
    return &DirectoryStore{messenger}
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

func (d *DirectoryStore) LoadTemplate(name, projectName, destinationPath string, fields map[string]string) (ok bool) {
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
        d.checkPath(sourcePath)
        d.checkPath(targetPath)

        return false
    default:
        d.messenger.Send(messaging.Message{
            Message: "An unforseen error has occurred",
            Severity: messaging.Error,
        })

        return false
    }

    return true
}

func (d *DirectoryStore) getSourcePath(name string) (string, bool) {
    executablePath, err := os.Executable()

    if err != nil {
        d.messenger.Send(messaging.Message{
            Message: "Cannot retrieve the executable path",
        })

        return "", false
    }

    return filepath.Join(filepath.Dir(executablePath), templatesFolderName, name), true
}

func (d *DirectoryStore) checkPath(path string) {
    _, err := os.Stat(path)

    if err != nil {
        d.messenger.Send(messaging.Message{
            Message: "The path '" + path + "' does not exist",
            Severity: messaging.Error,
        })
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
