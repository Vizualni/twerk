package twerk

// rename to twerk...twerk twerk twerk

import (
	"github.com/vizualni/twerk/callable"
	"github.com/vizualni/twerk/math"
	"log"
	"reflect"
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
	// Stop()
}

type twerk struct {
	callable *callable.Callable

	config Config

	jobListener chan jobInstruction

	liveWorkersNum      *atomicNumber
	currentlyWorkingNum *atomicNumber

	broadcastDie chan bool
}

type jobInstruction struct {
	arguments []reflect.Value
	returnTo  chan []interface{}
}

// Constructor for the Twerker.
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
	}

	twrkr.startInBackground()

	return twrkr, nil
}

// Starts another goroutine in background which monitors all
// others and has responsibility to scale up or down when necessary
func (twrkr *twerk) startInBackground() {

	twrkr.doINeedToStartMissingOnes()

	tick := time.NewTicker(twrkr.config.Refresh)

	go func() {
		defer tick.Stop()
		for range tick.C {

			twrkr.printStatus()

			if twrkr.doINeedToStartMissingOnes() {
				continue
			}

			if twrkr.doWeHaveTooLittleWorkers() {
				continue
			}

			if twrkr.doWeNeedToKillSomeWorkers() {
				continue
			}

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

	log.Printf("Live: %d; Working: %d; Idle: %d, Jobs in queue: %d; Min: %d; Max: %d",
		live, working, idle, inQueue, twrkr.config.Min, twrkr.config.Max)
}

// Checks if there are less than minimum number of workers available.
// If yes, it startes them.
func (twrkr *twerk) doINeedToStartMissingOnes() bool {
	live := twrkr.liveWorkersNum.Get()
	min := twrkr.config.Min

	if live < min {
		twrkr.startWorkers(min - live)
		return true
	}

	return false
}

// Checks if there are more jobs than workers. Scales workers up to the
// maximum available number.
func (twrkr *twerk) doWeHaveTooLittleWorkers() bool {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()
	inQueue := len(twrkr.jobListener)

	idle := live - working

	if idle >= inQueue {
		return false
	}

	howMuchToStart := math.Min(twrkr.config.Max-live, inQueue)

	twrkr.startWorkers(howMuchToStart)

	return true
}

// Are there way too many workers alive who are not doing anything?
// If yes, kill them!
func (twrkr *twerk) doWeNeedToKillSomeWorkers() bool {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()

	idle := live - working
	min := twrkr.config.Min

	if idle <= min {
		return false
	}

	twrkr.stopWorkers(idle - min)

	return true
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
		twrkr.liveWorkersNum.Incr()
		defer twrkr.liveWorkersNum.Decr()

		for {
			select {
			case job := <-twrkr.jobListener:
				twrkr.currentlyWorkingNum.Incr()
				returnValues := twrkr.callable.CallFunction(job.arguments)
				if len(returnValues) > 0 {
					go func() {
						job.returnTo <- returnValues
						close(job.returnTo)
					}()
				}
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

	argumentValues, err := twrkr.callable.TransformToValues(args...)

	if err != nil {
		return nil, err
	}

	var returnToChan chan []interface{}

	if twrkr.callable.NumberOfReturnValues() > 0 {
		returnToChan = make(chan []interface{})
	}

	newJobPass := jobInstruction{
		arguments: argumentValues,
		returnTo:  returnToChan,
	}

	go func() {
		twrkr.jobListener <- newJobPass
	}()

	return returnToChan, nil
}

// Waits until there are no more jobs in the queue, there are no more working workers and the alive number of workers
// is the minimum possible.
func (twrkr *twerk) Wait() {
	ticker := time.NewTicker(100 * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		if len(twrkr.jobListener) == 0 && twrkr.liveWorkersNum.Get() == twrkr.config.Min && twrkr.currentlyWorkingNum.Get() == 0 {
			return
		}
	}
}
