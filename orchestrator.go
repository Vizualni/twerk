package twerk

import (
	"github.com/vizualni/twerk/math"
)

type Orchestrator interface {
	Calculate(status Status) (startStopNum int, err error)
}

type defaultOrchestrator struct {
	config *Config
}

func (orch *defaultOrchestrator) Calculate(status Status) (startStopNum int, err error) {
	n := orch.doINeedToStartMissingOnes(status)
	if n != 0 {
		return n, nil
	}
	return -orch.doWeNeedToKillSomeWorkers(status), nil
}

// Checks if there are less than minimum number of workers available.
// If yes, it startes them.
func (orch *defaultOrchestrator) doINeedToStartMissingOnes(status Status) int {
	toStart := 0
	live := status.Live()
	max := orch.config.Max
	inQueue := math.Min(status.JobsInQueue(), max)

	if inQueue >= live {
		toStart += inQueue - live
	}

	return toStart
}

// Are there way too many workers alive who are not doing anything?
// If yes, kill them!
func (orch *defaultOrchestrator) doWeNeedToKillSomeWorkers(status Status) int {
	idle := status.Idle()
	inQueue := status.JobsInQueue()

	if idle <= inQueue {
		return 0
	}

	toKill := idle - inQueue

	return toKill
}
