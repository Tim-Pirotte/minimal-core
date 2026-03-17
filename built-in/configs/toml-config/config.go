package tomlconfig

import (
	"maps"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var configPath = filepath.Join(".", "commands")
var globalConfig = filepath.Join("global")

type TOMLConfig struct {
	defaultLocation string
	cache map[string]map[string]map[string]any
}

func NewConfig(defaultLocation string) *TOMLConfig {
	return &TOMLConfig{defaultLocation, nil}
}

func (t *TOMLConfig) Get(source, section, field string) (value any, ok bool) {
	if t.cache == nil {
		t.cache = map[string]map[string]map[string]any{}
	}

	if _, loaded := t.cache[source]; !loaded {
		if err := t.loadFile(source); err != nil {
			return nil, false
		}
	}

	sectionContent, ok := t.cache[source][section]
	
	if !ok {
		return nil, false
	}

	value, ok = sectionContent[field]

	if !ok {
		return nil, false
	}

	return value, true
}

func (t *TOMLConfig) loadFile(source string) error {
	global := make(map[string]map[string]any)

	globalPath := filepath.Join(configPath, globalConfig, source)

	if _, err := toml.DecodeFile(globalPath, &global); err != nil && !os.IsNotExist(err) {
		return err
	}

	local := make(map[string]map[string]any)

	localPath := filepath.Join(configPath, t.defaultLocation, source)

	if _, err := toml.DecodeFile(localPath, &local); err != nil {
		return err
	}

	for sectionName, fields := range local {
		if _, sectionExists := global[sectionName]; !sectionExists {
			global[sectionName] = map[string]any{}
		}

		maps.Copy(global[sectionName], fields)
	}

	t.cache[source] = global

	return nil
}
