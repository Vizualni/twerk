package twerk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestThatNonFunctionReturnsAnError(t *testing.T) {
	pool, err := New(123, DefaultConfig)

	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestThatInvalidConfigReturnsAnError(t *testing.T) {
	pool, err := New(func() {}, Config{})

	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestThatFunctionWithoutReturnValuesHasNilChannel(t *testing.T) {
	pool, err := New(func() {}, DefaultConfig)

	assert.NoError(t, err)
	assert.NotNil(t, pool)

	resChan, err := pool.Work()

	assert.NoError(t, err)
	assert.Nil(t, resChan)
}

func TestThatCallingFunctionWithIncorrectArgumentsReturnsAnError(t *testing.T) {
	pool, _ := New(func(a int) {}, DefaultConfig)

	resChan, err := pool.Work("this is not an int")

	assert.Error(t, err)
	assert.Nil(t, resChan)
}

func TestThatCallingWorkWithWrongNumberOfArgumentsReturnsAnError(t *testing.T) {
	pool, _ := New(func(a, b int) {}, DefaultConfig)

	resChan, err := pool.Work(1)

	assert.Error(t, err)
	assert.Nil(t, resChan)
}

func TestThatCallingAFunctionWithNoArgumentsWithArgumentsReturnsAnError(t *testing.T) {
	pool, _ := New(func() {}, DefaultConfig)

	resChan, err := pool.Work(1)

	assert.Error(t, err)
	assert.Nil(t, resChan)
}

func TestThatWorkerActuallyReturnsTheSameNumberOfReturnValues(t *testing.T) {
	pool, _ := New(func() (int, int, int) { return 1, 2, 3 }, DefaultConfig)

	resChan, err := pool.Work()

	assert.NoError(t, err)
	assert.NotNil(t, resChan)

	results := <-resChan
	assert.Len(t, results, 3)
	assert.Equal(t, 1, results[0])
	assert.Equal(t, 2, results[1])
	assert.Equal(t, 3, results[2])
}

func TestThatItWillScaleUpIfItCantWork(t *testing.T) {
	config := Config{}
	config.Min = 0
	config.Max = 1
	config.Refresh = 500 * time.Millisecond

	pool, _ := New(func(exit chan bool) {
		<-exit
	}, config)

	e1 := make(chan bool)
	e2 := make(chan bool)

	pool.Work(e1)

	didItWork := make(chan bool)

	go func(didItWork chan bool) {
		t1 := time.Now()
		pool.Work(e2)
		d := time.Now().Sub(t1)
		didItWork <- d > 450*time.Millisecond
	}(didItWork)

	time.Sleep(500 * time.Millisecond)
	e1 <- true
	e2 <- true

	assert.True(t, <-didItWork)
}
