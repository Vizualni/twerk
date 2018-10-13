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

func (an *atomicNumber) Incr() int {
	an.Lock()
	an.num++
	newValue := an.num
	an.Unlock()
	return newValue
}

func (an *atomicNumber) Decr() int {
	an.Lock()
	an.num--
	newValue := an.num
	an.Unlock()
	return newValue
}

func (an *atomicNumber) Get() int {
	an.Lock()
	value := an.num
	an.Unlock()
	return value
}
