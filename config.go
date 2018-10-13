package twerk

import (
	"fmt"
	"runtime"
	"time"
)

type Config struct {
	Max int
	Min int

	Refresh      time.Duration
	UseMyRefresh bool
}

var DefaultConfig = Config{
	Max: runtime.NumCPU(),
	Min: 0,

	Refresh: 1 * time.Second,
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

	if config.Refresh == 0 {
		return fmt.Errorf("set Refresh to some duration")
	}

	if config.Refresh < time.Millisecond && !config.UseMyRefresh {
		return fmt.Errorf("it is not the best to use that low refresh ratio. If you are sure then set UseMyRefresh to true")
	}

	return nil
}
