package config

import "github.com/BurntSushi/toml"

type Config interface {
	GetString(source, section, field string) (value any, ok bool)
}

func LoadConfig(configFile string, config any) error {
	_, err := toml.Decode(configFile, config)

	return err
}
