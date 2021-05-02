package retry

import (
	"fmt"
	"time"
)

const (
	MaxRetryForever    = 99999999
	MaxRetryThreeTimes = 3
	MaxRetryFiveTimes  = 5
	MaxRetryTenTimes   = 10
)

type AnyError struct {
	error
}

func DefaultAnyError() []error {
	return []error{AnyError{}}
}

type Config struct {
	MaxRetry  int
	DelayFunc func(attempt int) time.Duration
	ByErrors  []error
}

func DefaultConfig() Config {
	return Config{
		MaxRetry:  MaxRetryForever,
		DelayFunc: ConstantOneSecondDelayFunc,
		ByErrors:  DefaultAnyError(),
	}
}

// Validate called by Retry SetConfig to check every field is correct and valid
func (c *Config) Validate() error {

	if c.ByErrors == nil {
		// c.ByErrors = DefaultAnyError() <-- autocorrection ?
		return fmt.Errorf("use DefaultAnyError instead of nil")
	}

	if c.MaxRetry == 0 {
		return fmt.Errorf("MaxRetry must > 0")
	}

	if c.DelayFunc == nil {
		return fmt.Errorf("DelayFunc must not nil")
	}

	return nil
}
