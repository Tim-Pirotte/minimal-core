package pipeline

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestEmpty(t *testing.T) {
    p := NewPipeline()
    p.Run()
}

func TestRace(t *testing.T) {
    expectedCount := 100
    var count atomic.Uint64

    p := NewPipeline()

    stages := []Stage{
        {
            Name: "Stage 1",
            Units: []Unit{
                {
                    Name: "Unit 1",
                    Run: func() {
                        for i := range expectedCount {
                            p.AddUnitToCurrentStage(Unit{
                                Name: fmt.Sprintf("InnerUnit 1.%d", i),
                                Run: func() { count.Add(1) },
                            })
                        }
                    },
                },
                {
                    Name: "Unit 2",
                    Run: func() {
                        for i := range expectedCount {
                            p.AddUnitToCurrentStage(Unit{
                                Name: fmt.Sprintf("InnerUnit 2.%d", i),
                                Run: func() { count.Add(1) },
                            })
                        }
                    },
                },
            },
        },
        {
            Name: "Intermediate check",
            Units: []Unit{
                {
                    Name: "CheckCount",
                    Run: func() {
                        actualCount := count.Load()

                        if actualCount != 2 * uint64(expectedCount) {
                            t.Errorf("Expected %d but got %d", expectedCount, actualCount)
                        }
                    },
                },
            },
        },
        {
            Name: "Stage 2",
            Units: []Unit{
                {
                    Name: "Unit 3",
                    Run: func() {
                        for i := range expectedCount {
                            p.AddUnitToCurrentStage(Unit{
                                Name: fmt.Sprintf("InnerUnit 3.%d", i),
                                Run: func() { count.Add(1) },
                            })
                        }
                    },
                },
                {
                    Name: "Unit 4",
                    Run: func() {
                        for i := range expectedCount {
                            p.AddUnitToCurrentStage(Unit{
                                Name: fmt.Sprintf("InnerUnit 4.%d", i),
                                Run: func() { count.Add(1) },
                            })
                        }
                    },
                },
            },
        },
    }

    p.Stages = stages
    p.Run()

    actualCount := count.Load()

    if actualCount != 4 * uint64(expectedCount) {
        t.Errorf("Expected %d but got %d", expectedCount, actualCount)
    }
}
