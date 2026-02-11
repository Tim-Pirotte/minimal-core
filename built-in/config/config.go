package config

import "github.com/BurntSushi/toml"

func LoadConfig(configFile string, config any) error {
	_, err := toml.Decode(configFile, config)

	return err
}
