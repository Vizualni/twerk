package callable_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vizualni/twerk/callable"
	"testing"
)

func TestNewWithVoidFunc(t *testing.T) {
	call, err := callable.New(func() {})

	assert.Nil(t, err)
	assert.NotNil(t, call)
	assert.Equal(t, 0, call.NumberOfReturnValues())
}

func TestNewWithNotAFunc(t *testing.T) {
	call, err := callable.New(123)

	assert.Error(t, err)
	assert.Nil(t, call)
}

func TestNewWithLittleComplexFunction(t *testing.T) {
	call, err := complexFunction()

	assert.Nil(t, err)
	assert.NotNil(t, call)
	assert.Equal(t, 2, call.NumberOfReturnValues())
}

func TestCallWithIncorrectArgumentsReturnsAnError(t *testing.T) {
	call, _ := complexFunction()

	_, err := callFunc(call, 123)
	assert.Error(t, err)

	_, err = callFunc(call, 123, 456, "bla")
	assert.Error(t, err)

	_, err = callFunc(call, 123, 456, "bla", &testStruct{}) // note that it's an address
	assert.Error(t, err)
}

func TestCallWithCorrectArgumentsDoesNotReturnAnError(t *testing.T) {
	call, _ := complexFunction()

	results, err := callFunc(call, 123, 456, "bla", testStruct{})
	assert.Nil(t, err)
	assert.NotNil(t, results)
}

func TestCallWithCorrectArgumentsReturnsCorrectResult(t *testing.T) {
	call, _ := complexFunction()

	results, _ := callFunc(call, 12, 4, "bla", testStruct{})
	assert.Len(t, results, 2)
	assert.Equal(t, 16, results[0])
	assert.Equal(t, fmt.Errorf("bla"), results[1])
}

type testStruct struct {
	a, b int
}

func complexFunction() (*callable.Callable, error) {
	return callable.New(func(a, b int, c string, d testStruct) (int, error) {
		return a + b, fmt.Errorf(c)
	})
}

func callFunc(call *callable.Callable, arguments ...interface{}) ([]interface{}, error) {
	values, err := call.TransformToValues(arguments...)

	if err != nil {
		return nil, err
	}

	return call.CallFunction(values), nil
}
