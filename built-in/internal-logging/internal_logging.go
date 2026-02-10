package logging

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
)

var rootLogger zerolog.Logger
var activeBuffer *ringBuffer
var initCalled bool

type SourceGenerator struct {
    path            []string
    declaredSources map[string]int
    authentic       bool
}

func Init(maxBytes uint) SourceGenerator {
    if initCalled {
        panic("logging.Init has already been called")
    }

	activeBuffer = newRingBuffer(maxBytes)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    rootLogger = zerolog.New(activeBuffer).With().Timestamp().Logger()

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
        s.declaredSources[name] = 1
    }

    newPath := append([]string(nil), s.path...)
    newPath = append(newPath, name)

    return rootLogger.With().Strs("source", newPath).Logger(), SourceGenerator{newPath, make(map[string]int), true}
}

func WriteTo(w io.Writer) (int64, error) {
    if activeBuffer == nil {
        panic("the ring buffer for the logger has not been initialized")
    }

    return activeBuffer.writeTo(w)
}
