package twerk

// rename to twerk...twerk twerk twerk

import (
	"github.com/vizualni/twerk/callable"
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

	for i := 0; i < twrkr.config.Min; i++ {
		twrkr.startWorker()
	}

	tick := time.NewTicker(twrkr.config.Refresh)

	go func() {
		defer tick.Stop()
		for range tick.C {

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

	remainingInQeueu := inQueue - live

	if remainingInQeueu == 0 {
		return false
	}

	howMuchToStart := remainingInQeueu

	if howMuchToStart > twrkr.config.Max {
		howMuchToStart = twrkr.config.Max
	}

	twrkr.startWorkers(howMuchToStart)

	return true
}

func (twrkr *twerk) doWeNeedToKillSomeWorkers() bool {
	live := twrkr.liveWorkersNum.Get()
	working := twrkr.currentlyWorkingNum.Get()
	inQueue := len(twrkr.jobListener)

	idle := live - working

	if idle < inQueue {
		return false
	}

	twrkr.stopWorkers(inQueue - idle)

	return true
}

func (twrkr *twerk) startWorkers(n int) {
	if n <= 0 {
		return
	}
	for i := 0; i < n; i++ {
		twrkr.startWorker()
	}
}

func (twrkr *twerk) stopWorkers(n int) {
	if n <= 0 {
		return
	}
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
