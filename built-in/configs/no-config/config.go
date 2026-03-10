package noconfig

import "errors"

var ErrNoConfig = errors.New("no config loader")

type NoConfig struct{}

func NewNoConfig() *NoConfig {
	return &NoConfig{}
}

func (*NoConfig) Get(any, string) error {
	return ErrNoConfig
}
