package logging

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
)

var rootLogger zerolog.Logger
var initCalled bool

type SourceGenerator struct {
    path            []string
    declaredSources map[string]int
    authentic       bool
}

func Init(target io.Writer) SourceGenerator {
    if initCalled {
        panic("logging.Init has already been called")
    }

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    rootLogger = zerolog.New(target).With().Timestamp().Logger()

    // The authentic field prevents init working multiple times 
    // and external packages creating SourceGenerators themselves
    sourceGenerator := SourceGenerator{[]string{}, make(map[string]int), !initCalled}

    initCalled = true

    return sourceGenerator
}

func (s *SourceGenerator) GetLogger(name string) (zerolog.Logger, SourceGenerator) {
    if !s.authentic {
        panic("this is not a SourceGenerator given by the logging package")
    }

    name = strings.TrimSpace(name)
    
    if name == "" {
        name = "unnamed"
    }

    if count, ok := s.declaredSources[name]; ok {
        s.declaredSources[name] = count + 1
        name = fmt.Sprintf("%s#%d", name, count+1)
    } else {
        s.declaredSources[name] = 0
    }

    newPath := append([]string(nil), s.path...)
    newPath = append(newPath, name)

    return rootLogger.With().Strs("source", newPath).Logger(), SourceGenerator{newPath, make(map[string]int), true}
}
