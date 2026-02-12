package templating

import (
	"os"
	"path/filepath"
	"runtime"
)

const appName = "minimal"

var dataDirectories = map[string]string{
	"linux": getEnvOrFallback("XDG_DATA_HOME", filepath.Join(os.Getenv("HOME"), ".local", "share")),
	"windows": os.Getenv("LocalAppData"),
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support"),
}

var fallbackDataDirectory = getEnvOrFallback("XDG_DATA_HOME", filepath.Join(os.Getenv("HOME"), ".local", "share"))

var configDirectories = map[string]string{
	"linux": getEnvOrFallback("XDG_CONFIG_HOME", filepath.Join(os.Getenv("HOME"), ".config")),
	"windows": os.Getenv("AppData"),
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Preferences"),
}

var fallbackConfigDirectory = getEnvOrFallback("XDG_CONFIG_HOME", filepath.Join(os.Getenv("HOME"), ".config"))

func getEnvOrFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func GetDataDirectory() string {
	path, ok := dataDirectories[runtime.GOOS]

	if !ok {
		path = fallbackDataDirectory
	}

	return filepath.Join(path, appName)
}

func GetConfigDirectory() string {
	path, ok := configDirectories[runtime.GOOS]

	if !ok {
		path = fallbackConfigDirectory
	}

	return filepath.Join(path, appName)
}
