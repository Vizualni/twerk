package twerk

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
