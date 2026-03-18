package noconfig

type NoConfig struct{}

func NewNoConfig() *NoConfig {
	return &NoConfig{}
}

func (*NoConfig) SetLocalConfigSource(location string) {}

func (*NoConfig) Get(source, section, field string) (value any, ok bool) {
	return nil, false
}
