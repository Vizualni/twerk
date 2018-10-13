package twerk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitNumber(t *testing.T) {
	num := newAtomicNumber(123)
	assert.Equal(t, 123, num.Get())
}

func TestThatIncreaseActuallyIncreasesANumber(t *testing.T) {
	num := newAtomicNumber(123)
	num.Incr()
	assert.Equal(t, 124, num.Get())
}

func TestThatDecreaseActuallyDecreasesANumber(t *testing.T) {
	num := newAtomicNumber(123)
	num.Decr()
	assert.Equal(t, 122, num.Get())
}

func TestThatMultipleGoroutinesAccessingItCalculatesItCorrectly(t *testing.T) {
	num := newAtomicNumber(0)

	increaseFunc := func(number *atomicNumber) {
		number.Incr()
	}

	go increaseFunc(num)
	go increaseFunc(num)

	// not the best way but I don't want to introduce locking here
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 2, num.Get())
}

func TestThatMultipleGoroutinesAccessingItCalculatesItCorrectlyDecreasing(t *testing.T) {
	num := newAtomicNumber(0)

	decreaseFunc := func(number *atomicNumber) {
		number.Decr()
	}

	go decreaseFunc(num)
	go decreaseFunc(num)

	// not the best way but I don't want to introduce locking here
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, -2, num.Get())
}
