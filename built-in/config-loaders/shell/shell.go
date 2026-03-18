package shell

import configloader "minimal/minimal-core/built-in/config-loader"

type Shell struct {
	Implementation configloader.ConfigLoader
}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) SetLocalConfigSource(location string) {
	s.SetLocalConfigSource(location)
}

func (s *Shell) Get(source, section, field string) (value any, ok bool) {
	return s.Get(source, section, field)
}
