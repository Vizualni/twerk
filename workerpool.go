package twerk

// rename to twerk...twerk twerk twerk

import (
	"fmt"
	"github.com/vizualni/twerk/callable"
	"log"
	"reflect"
	"sync"
	"time"
)

// Twerker (or Worker) interface.
type Twerker interface {

	// Calls the actual function with the given arguments.
	Work(args ...interface{}) (<-chan []interface{}, error)

	// Waits until there are no more jobs in the queue, then it releases.
	Wait()

	// Stops accepting new jobs and finishes everything.
	// Work cannot be called after Stop has already been called.
	Stop()
}

type twerk struct {
	callable *callable.Callable

	config Config

	jobListener chan jobInstruction

	liveWorkersNum      *atomicNumber
	currentlyWorkingNum *atomicNumber

	broadcastDie chan bool

	stop bool

	orchestrator Orchestrator

	stopTheWorldLock *sync.RWMutex
}

type jobInstruction struct {
	arguments []reflect.Value
	returnTo  chan []interface{}
}

// New is the constructor for the Twerker.
func New(v interface{}, config Config) (*twerk, error) {

	callableFunc, err := callable.New(v)

	if err != nil {
		return nil, err
	}

	err = isValid(config)

	if err != nil {
		return nil, err
	}

	twrkr := &twerk{
		callable: callableFunc,

		config: config,

		jobListener: make(chan jobInstruction, config.Max),

		liveWorkersNum:      newAtomicNumber(0),
		currentlyWorkingNum: newAtomicNumber(0),

		broadcastDie: make(chan bool),

		orchestrator: &defaultOrchestrator{config: &config},

		stopTheWorldLock: &sync.RWMutex{},
	}

	twrkr.startInBackground()

	return twrkr, nil
}

// Starts another goroutine in background which monitors all
// others and has responsibility to scale up or down when necessary
func (twrkr *twerk) startInBackground() {

	twrkr.changeCapacity()

	tick := time.NewTicker(twrkr.config.Refresh)

	go func() {
		defer tick.Stop()
		for range tick.C {
			if twrkr.stop {
				return
			}
			twrkr.printStatus()
			twrkr.changeCapacity()
		}
	}()
}

// This is just a debug line.
// This should be opt-in from the config.
func (twrkr *twerk) printStatus() {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()
	inQueue := len(twrkr.jobListener)
	idle := live - working

	log.Printf("Live: %d; Working: %d; Idle: %d, Jobs in queue: %d; Max: %d",
		live, working, idle, inQueue, twrkr.config.Max)
}

// This is just a debug line.
// This should be opt-in from the config.
func (twrkr *twerk) status() Status {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()
	inQueue := len(twrkr.jobListener)

	return Status{
		live:        live,
		working:     working,
		jobsInQueue: inQueue,
	}
}

// Starts n workers.
func (twrkr *twerk) startWorkers(n int) {
	if n <= 0 {
		return
	}
	log.Printf("Starting %d workers\n", n)
	for i := 0; i < n; i++ {
		twrkr.startWorker()
	}
}

// Stops n workers.
func (twrkr *twerk) stopWorkers(n int) {
	if n <= 0 {
		return
	}
	log.Printf("Stopping %d workers\n", n)
	for i := 0; i < n; i++ {
		twrkr.stopWorker()
	}
}

// Starts a single worker.
// Also increases number of live workers when started.
// When job starts processing then it increases number of currently working workers.
func (twrkr *twerk) startWorker() {

	go func() {
		twrkr.waitOnWorld()
		twrkr.liveWorkersNum.Incr()
		defer func() {
			twrkr.waitOnWorld()
			twrkr.liveWorkersNum.Decr()
		}()
		for {
			select {
			case job, _ := <-twrkr.jobListener:
				twrkr.waitOnWorld()
				twrkr.currentlyWorkingNum.Incr()
				returnValues := twrkr.callable.CallFunction(job.arguments)
				if len(returnValues) > 0 {
					go func() {
						job.returnTo <- returnValues
						close(job.returnTo)
					}()
				}
				twrkr.waitOnWorld()
				twrkr.currentlyWorkingNum.Decr()
			case <-twrkr.broadcastDie:
				// somebody requested that we die
				return
			}
		}
	}()

}

// Stops worker, or if there arent any available,
// does nothing ofter one millisecond.
// Not the best approach
func (twrkr *twerk) stopWorker() {
	select {
	case twrkr.broadcastDie <- true:
		// nice
	case <-time.After(time.Millisecond):
		// nobody is listening
	}
}

// The actual work function that calls the function with given arguments
func (twrkr *twerk) Work(args ...interface{}) (<-chan []interface{}, error) {

	if twrkr.stop {
		return nil, fmt.Errorf("twerker has been stopped. no more work can be done with this twerker")
	}

	argumentValues, err := twrkr.callable.TransformToValues(args...)

	if err != nil {
		return nil, err
	}

	var returnToChan chan []interface{}

	if twrkr.callable.NumberOfReturnValues() > 0 {
		returnToChan = make(chan []interface{})
	}

	newJobInstruction := jobInstruction{
		arguments: argumentValues,
		returnTo:  returnToChan,
	}

	go func() {
		twrkr.jobListener <- newJobInstruction
	}()

	return returnToChan, nil
}

// Waits until there are no more jobs in the queue, there are no more working workers and the alive number of workers
// is the minimum possible.
func (twrkr *twerk) Wait() {
	if twrkr.stop {
		return
	}
	ticker := time.NewTicker(100 * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		if len(twrkr.jobListener) == 0 && twrkr.liveWorkersNum.Get() == 0 && twrkr.currentlyWorkingNum.Get() == 0 {
			return
		}
	}
}

// Stops
func (twrkr *twerk) Stop() {
	twrkr.Wait()
	twrkr.stop = true // this should be included in stop the world
	twrkr.stopWorkers(twrkr.liveWorkersNum.Get())
	close(twrkr.jobListener)
}

func (twrkr *twerk) changeCapacity() {
	// stop the world
	twrkr.stopTheWorldLock.Lock()
	defer twrkr.stopTheWorldLock.Unlock()

	status := twrkr.status()
	addRemove, _ := twrkr.orchestrator.Calculate(status)

	switch {
	case addRemove > 0:
		twrkr.startWorkers(addRemove)
	case addRemove < 0:
		twrkr.stopWorkers(-addRemove)
	}
}

func (twrkr *twerk) waitOnWorld() {
	twrkr.stopTheWorldLock.RLock()
	defer twrkr.stopTheWorldLock.RUnlock()
}
