package sources

import "math"

type Sources struct {
	Sources []Source
}

type Source struct {
	Name    string
	Content string
}

func New() *Sources {
	return &Sources{[]Source{}}
}

func (s *Sources) Add(source Source) bool {
	if len(source.Content) > math.MaxUint32 {
		return false
	}

	s.Sources = append(s.Sources, source)

	return true
}
