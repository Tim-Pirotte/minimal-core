package tomlconfig

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var configPath = filepath.Join(".", "commands")
var globalConfig = filepath.Join("global")

type TOMLConfig struct {
	defaultLocation string
}

func NewConfig(defaultLocation string) *TOMLConfig {
	return &TOMLConfig{}
}

func (t *TOMLConfig) Get(loadInto any, source string) error {
	_, err := toml.Decode(filepath.Join(configPath, globalConfig, source), loadInto)

	if err != nil {
		return err
	}

	_, err = toml.Decode(filepath.Join(configPath, t.defaultLocation, source), loadInto)

	return err
}
