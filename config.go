package twerk

import (
	"fmt"
	"runtime"
	"time"
)

// Config settings for the Twerker
type Config struct {

	// Defines a maximum number of workers which are allowed to work at the same time
	Max int

	// After each Refresh duration, it will check if it needs to scale up or down
	Refresh time.Duration
	// Set it to true only if you want to have really low value for the Refresh interval
	UseMyRefresh bool

	// If set to true, twerker will periodically print log values
	Debug bool
}

// DefaultConfig is configuration that you can use instead of creating your own.
// Maximum is defined as number of CPU cores you have.
var DefaultConfig = Config{
	Max: runtime.NumCPU(),

	Refresh: 1 * time.Second,
}

func isValid(config Config) error {

	if config.Max <= 0 {
		return fmt.Errorf("max (%d) must be greater than zero", config.Max)
	}

	if config.Refresh == 0 {
		return fmt.Errorf("set Refresh to some duration")
	}

	if config.Refresh < time.Millisecond && !config.UseMyRefresh {
		return fmt.Errorf("it is not the best to use that low refresh ratio. If you are sure then set UseMyRefresh to true")
	}

	return nil
}
