package configloader

// TODO remove
type ConfigLoader interface {
	SetLocalConfigSource(location string)
	Get(source, section, field string) (value any, ok bool)
}
