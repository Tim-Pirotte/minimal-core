package shell

import configloader "minimal/minimal-core/built-in/config-loader"

type ShellConfigLoader struct {
	Implementation configloader.ConfigLoader
}

func NewShellConfigLoader(defaultImplementation configloader.ConfigLoader) *ShellConfigLoader {
	return &ShellConfigLoader{defaultImplementation}
}

func (s *ShellConfigLoader) SetLocalConfigSource(location string) {
	s.Implementation.SetLocalConfigSource(location)
}

func (s *ShellConfigLoader) Get(source, section, field string) (value any, ok bool) {
	return s.Implementation.Get(source, section, field)
}
