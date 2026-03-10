package config

import "github.com/BurntSushi/toml"

type Config interface {
	Get(key, source string) (any, error)
}

func LoadConfig(configFile string, config any) error {
	_, err := toml.Decode(configFile, config)

	return err
}
