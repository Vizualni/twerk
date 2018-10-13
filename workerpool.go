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

	liveWorkersNum *atomicNumber

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

		jobListener: make(chan jobInstruction, 15),

		liveWorkersNum: newAtomicNumber(0),

		broadcastDie: make(chan bool),
	}

	twrkr.startInBackground()

	return twrkr, nil
}

func (twrkr *twerk) startInBackground() {

	for i := 0; i < twrkr.config.Min; i++ {
		twrkr.startWorker()
	}

	tick := time.NewTicker(1000 * time.Millisecond)

	go func() {
		defer tick.Stop()
		for range tick.C {

			inQueue := len(twrkr.jobListener)
			live := twrkr.liveWorkersNum.Get()
			min := twrkr.config.Min

			// some of the workers died?
			// better startInBackground them again
			if live < min {
				for i := live; i < twrkr.config.Min; i++ {
					twrkr.startWorker()
				}
				continue
			}

			// we are far behind and need to startInBackground new workers
			if inQueue >= live {
				howMuchWorkersToStart := inQueue - live
				for i := 0; i < howMuchWorkersToStart; i++ {
					twrkr.startWorker()
				}
				continue
			}

			killThisMany := live - inQueue - min

			if killThisMany > 0 {
				for i := 0; i < killThisMany; i++ {
					twrkr.killWorker()
				}
			}

		}
	}()
}

func (twrkr *twerk) startWorker() {

	go func() {
		twrkr.liveWorkersNum.Incr()
		defer twrkr.liveWorkersNum.Decr()

		for {
			select {
			case job := <-twrkr.jobListener:

				returnValues := twrkr.callable.CallFunction(job.arguments)

				if len(returnValues) > 0 {
					go func() {
						// check if channel is closed
						_, closed := <-job.returnTo

						if !closed {
							job.returnTo <- returnValues
							close(job.returnTo)
						}
					}()
				}
			case <-twrkr.broadcastDie:
				// somebody requested that we die
				return
			}
		}
	}()

}

func (twrkr *twerk) killWorker() {
	twrkr.broadcastDie <- true
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
