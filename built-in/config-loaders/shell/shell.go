package shell

import configloader "minimal/minimal-core/built-in/config-loader"

type Shell struct {
	Implementation configloader.ConfigLoader
}

func NewShell(defaultImplementation configloader.ConfigLoader) *Shell {
	return &Shell{defaultImplementation}
}

func (s *Shell) SetLocalConfigSource(location string) {
	s.Implementation.SetLocalConfigSource(location)
}

func (s *Shell) Get(source, section, field string) (value any, ok bool) {
	return s.Implementation.Get(source, section, field)
}
