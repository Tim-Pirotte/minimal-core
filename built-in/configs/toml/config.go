package toml

import (
	"maps"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

var configPath = filepath.Join(".", "commands")
var globalConfig = filepath.Join("global")

type TOMLConfig struct {
	location string
	cache map[string]map[string]map[string]any
	logger zerolog.Logger
}

func NewConfigLoader(srcGen *logging.SourceGenerator) *TOMLConfig {
	logger, _ := srcGen.GetLogger("tomlConfigLoader")
	
	return &TOMLConfig{"", nil, logger}
}

func (t *TOMLConfig) SetLocalConfigSource(location string) {
	t.location = location
	t.logger.Debug().Msg("new config location set")
}

func (t *TOMLConfig) Get(source, section, field string) (value any, ok bool) {
	if t.cache == nil {
		t.cache = map[string]map[string]map[string]any{}
	}

	if _, loaded := t.cache[source]; !loaded {
		if err := t.loadFile(source); err != nil {
			t.logger.Error().
				Err(err).
				Str("source", source).
				Str("section", section).
				Str("field", field).
				Msg("error loading toml file")

			return nil, false
		}
	}

	sectionContent, ok := t.cache[source][section]
	
	if !ok {
		t.logger.Error().
			Str("source", source).
			Str("section", section).
			Str("field", field).
			Msg("section not in config")
		
		return nil, false
	}

	value, ok = sectionContent[field]

	if !ok {
		t.logger.Error().
			Str("source", source).
			Str("section", section).
			Str("field", field).
			Msg("section not in config")
	
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

	localPath := filepath.Join(configPath, t.location, source)

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
