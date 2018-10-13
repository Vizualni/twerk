package twerk

import (
	"fmt"
	"runtime"
)

type Config struct {
	Max int
	Min int
}

var DefaultConfig = Config{
	Max: runtime.NumCPU(),
	Min: 0,
}

func isValid(config Config) error {

	if config.Max <= 0 {
		return fmt.Errorf("max (%d) must be greater than zero", config.Max)
	}
	if config.Min < 0 {
		return fmt.Errorf("min (%d) must be greater than or equal to zero", config.Min)
	}

	if config.Min > config.Max {
		return fmt.Errorf("it would be cool, but min (%d) can't be bigger than max (%d)", config.Min, config.Max)
	}

	return nil
}
