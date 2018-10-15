package twerk

import "sync"

type atomicNumber struct {
	sync.Mutex
	num int
}

func newAtomicNumber(num int) *atomicNumber {
	return &atomicNumber{
		num:   num,
		Mutex: sync.Mutex{},
	}
}

// Atomic increase
func (an *atomicNumber) Incr() int {
	an.Lock()
	an.num++
	newValue := an.num
	an.Unlock()
	return newValue
}

// Atomic decrease
func (an *atomicNumber) Decr() int {
	an.Lock()
	an.num--
	newValue := an.num
	an.Unlock()
	return newValue
}

// Gets number
func (an *atomicNumber) Get() int {
	an.Lock()
	value := an.num
	an.Unlock()
	return value
}
