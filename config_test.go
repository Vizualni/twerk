package twerk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatCorrectConfigDoesNotReturnAnError(t *testing.T) {
	config := Config{
		Min: 1,
		Max: 2,
	}

	assert.Nil(t, isValid(config))
}

func TestIncorrectConfigsReturnAnError(t *testing.T) {
	for _, config := range incorrectConfigs() {
		t.Run("incorrect_config_test", func(t *testing.T) {
			assert.Error(t, isValid(config), fmt.Sprintf("%+v", config))
		})
	}
}

func incorrectConfigs() []Config {
	return []Config{
		{
			Max: 0,
			Min: 0,
		},
		{
			Max: -1,
			Min: -1,
		},
		{
			Max: -1,
			Min: 0,
		},
		{
			Max: 0,
			Min: -1,
		},
		{
			Max: 10,
			Min: 20,
		},
		{
			Max: 1,
			Min: -1,
		},
		{
			Max: -1,
			Min: 1,
		},
	}
}
