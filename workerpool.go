package twerk

// rename to twerk...twerk twerk twerk

import (
	"github.com/vizualni/twerk/callable"
	"github.com/vizualni/twerk/math"
	"log"
	"reflect"
	"time"
)

type Twerker interface {
	Work(args ...interface{}) (<-chan []interface{}, error)
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

func (twrkr *twerk) printStatus() {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()
	inQueue := len(twrkr.jobListener)
	idle := live - working

	log.Printf("Live: %d; Working: %d; Idle: %d, Jobs in queue: %d; Min: %d; Max: %d",
		live, working, idle, inQueue, twrkr.config.Min, twrkr.config.Max)
}

func (twrkr *twerk) doINeedToStartMissingOnes() bool {
	live := twrkr.liveWorkersNum.Get()
	min := twrkr.config.Min

	if live < min {
		twrkr.startWorkers(min - live)
		return true
	}

	return false
}

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

func (twrkr *twerk) startWorkers(n int) {
	if n <= 0 {
		return
	}
	log.Printf("Starting %d workers\n", n)
	for i := 0; i < n; i++ {
		twrkr.startWorker()
	}
}

func (twrkr *twerk) stopWorkers(n int) {
	if n <= 0 {
		return
	}
	log.Printf("Stopping %d workers\n", n)
	for i := 0; i < n; i++ {
		twrkr.stopWorker()
	}
}

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

func (twrkr *twerk) stopWorker() {
	select {
	case twrkr.broadcastDie <- true:
		// nice
	case <-time.After(time.Millisecond):
		// nobody is listening
	}
}

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

	// blocks until it can write to a channel
	twrkr.jobListener <- newJobPass

	return returnToChan, nil
}
