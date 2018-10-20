package twerk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCorrectAndIncorrectConfigs(t *testing.T) {
	for _, config := range incorrectConfigs() {
		t.Run("incorrect_config_test", func(t *testing.T) {
			assert.Error(t, isValid(config), fmt.Sprintf("%+v", config))
		})
	}
	for _, config := range correctConfigs() {
		t.Run("correct_config_test", func(t *testing.T) {
			assert.NoError(t, isValid(config), fmt.Sprintf("%+v", config))
		})
	}
}

func incorrectConfigs() []Config {
	return []Config{
		{
			Max: -1,
		},
		{
			Max: 0,
		},
		{
			Max: 10,
		},
		{
			Max: 1,
		},
		{
			Max:          10,
			Refresh:      time.Nanosecond,
			UseMyRefresh: false,
		},
	}
}

func correctConfigs() []Config {
	return []Config{
		{
			Max:     2,
			Refresh: time.Second,
		},
		{
			Max:          2,
			Refresh:      time.Nanosecond,
			UseMyRefresh: true,
		},
	}
}
