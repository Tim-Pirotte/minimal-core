package pipeline

import (
	"fmt"
	"sync"
)

type Pipeline struct {
    stages []Stage
    wg     sync.WaitGroup
}

type Stage struct {
    Name  string
    Units []Unit
}

type Unit struct {
    Name string
    Run func()
}

func NewPipeline() *Pipeline {
    return &Pipeline{}
}

func (p *Pipeline) Run() {
    for _, stage := range p.stages {
        for _, unit := range stage.Units {
            p.wg.Go(unit.Run)
        }

        p.wg.Wait()
    }
}

// Should only be called from within a running unit
func (p *Pipeline) AddUnitToCurrentStage(unit Unit) {
    p.wg.Go(unit.Run)
}

func (s *Stage) logStartStage() {
    // TODO use proper messaging
    fmt.Printf("Starting stage: %s\n", s.Name)
}

func (u *Unit) logStartUnit() {
    // TODO use proper messaging
    fmt.Printf("Starting unit: %s\n", u.Name)
}
