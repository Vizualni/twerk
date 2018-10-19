package twerk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type configStatusPair struct {
	config Config
	status Status

	expectedNum int
}

func TestTHatCorrectNumbersAreReturned(t *testing.T) {
	//t.Parallel()
	for _, pair := range getTestData() {
		config := pair.config
		status := pair.status
		expected := pair.expectedNum
		t.Run("calculate-correct", func(t *testing.T) {
			do := defaultOrchestrator{config: &config}
			res, _ := do.Calculate(status)
			assert.Equal(t, expected, res, "Config: %+v, Status: %+v", config, status)
		})
	}

}

func getTestData() []configStatusPair {
	return []configStatusPair{
		{
			config:      newConfigMinMax(0, 1),
			status:      newStatus(1, 0, 0),
			expectedNum: -1,
		},
		{
			config:      newConfigMinMax(2, 2),
			status:      newStatus(1, 0, 0),
			expectedNum: 1,
		},
		{
			config:      newConfigMinMax(2, 6),
			status:      newStatus(4, 1, 10),
			expectedNum: 2,
		},
		{
			config:      newConfigMinMax(2, 15),
			status:      newStatus(2, 1, 10),
			expectedNum: 8,
		},
		{
			config:      newConfigMinMax(2, 5),
			status:      newStatus(5, 2, 0),
			expectedNum: -3,
		},
		{
			config:      newConfigMinMax(2, 10),
			status:      newStatus(5, 2, 1),
			expectedNum: -2,
		},
		{
			config:      newConfigMinMax(0, 10),
			status:      newStatus(1, 0, 0),
			expectedNum: -1,
		},
		{
			config:      newConfigMinMax(0, 10),
			status:      newStatus(1, 1, 0),
			expectedNum: 0,
		},
		{
			config:      newConfigMinMax(5, 10),
			status:      newStatus(5, 0, 3),
			expectedNum: 0,
		},
		{
			config:      newConfigMinMax(5, 10),
			status:      newStatus(5, 0, 5),
			expectedNum: 0,
		},
		{
			config:      newConfigMinMax(5, 10),
			status:      newStatus(5, 0, 6),
			expectedNum: 1,
		},
		{
			config:      newConfigMinMax(5, 10),
			status:      newStatus(5, 5, 99),
			expectedNum: 5,
		},
	}
}

func newConfigMinMax(min, max int) Config {
	return Config{
		Min: min,
		Max: max,
	}
}
func newStatus(live, working, jobsInQueue int) Status {
	return Status{
		live:        live,
		working:     working,
		jobsInQueue: jobsInQueue,
	}
}
