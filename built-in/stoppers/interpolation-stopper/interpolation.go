package interpolation

import "minimal/minimal-core/built-in/lexer"

type InterpolationStopper struct{
    nesting uint
}

func NewInterpolationStopper() *InterpolationStopper {
	return &InterpolationStopper{1}
}

func (i *InterpolationStopper) End(s *lexer.LexerJob) bool {
	if s.Position >= uint(len(s.Data)) {
        return true
    }

    switch s.Data[s.Position] {
    case '{':
        i.nesting++
    case '}':
        i.nesting--
    }

    return i.nesting == 0
}
