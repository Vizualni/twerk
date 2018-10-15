package callable

import (
	"fmt"
	"reflect"
)

type Callable struct {
	callable reflect.Value

	argumentTypes []reflect.Type

	numberOfReturnValues int
}

func New(v interface{}) (*Callable, error) {

	reflectValue := reflect.ValueOf(v)

	if reflectValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function")
	}

	return &Callable{
		callable:             reflectValue,
		argumentTypes:        funcArgumentTypes(reflectValue),
		numberOfReturnValues: numberOfReturnValues(reflectValue),
	}, nil
}

func (c *Callable) NumberOfReturnValues() int {
	return c.numberOfReturnValues
}

func (c *Callable) CallFunction(argumentValues []reflect.Value) (resultInterfaces []interface{}) {

	shouldReturn := c.numberOfReturnValues > 0

	results := c.callable.Call(argumentValues)
	if shouldReturn {
		for _, result := range results {
			resultInterfaces = append(resultInterfaces, result.Interface())
		}
	}

	return
}

func (c *Callable) TransformToValues(arguments ...interface{}) ([]reflect.Value, error) {

	if len(arguments) != len(c.argumentTypes) {
		return nil, fmt.Errorf("expected %d number of arguments, got %d", len(c.argumentTypes), len(arguments))
	}

	argumentValues := argumentsToValues(arguments...)

	// check if arguments match
	for i := range c.argumentTypes {

		// todo ovo jos provjeri
		if c.argumentTypes[i].Kind() == reflect.Interface {
			if !argumentValues[i].Type().Implements(c.argumentTypes[i]) {
				return nil, fmt.Errorf("%s argument does not implement %s", argumentValues[i].String(), c.argumentTypes[i].String())
			}
			continue
		}
		if c.argumentTypes[i].Kind() != argumentValues[i].Kind() {
			return nil, fmt.Errorf(
				"argument with index %d of type %s does not match with expected type %s",
				i, c.argumentTypes[i].Kind().String(), argumentValues[i].Kind().String(),
			)
		}
	}

	return argumentValues, nil
}

func funcArgumentTypes(reflectValue reflect.Value) (in []reflect.Type) {

	reflectType := reflectValue.Type()
	numArguments := reflectType.NumIn()

	for i := 0; i < numArguments; i++ {
		in = append(in, reflectType.In(i))
	}

	return
}

func numberOfReturnValues(reflectValue reflect.Value) int {
	return reflectValue.Type().NumOut()
}

func argumentsToValues(args ...interface{}) (values []reflect.Value) {

	numArguments := len(args)

	for i := 0; i < numArguments; i++ {
		values = append(values, reflect.ValueOf(args[i]))
	}

	return
}
