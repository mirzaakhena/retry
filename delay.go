package retry

import (
	"math"
	"time"
)

func ConstantOneSecondDelayFunc(attempt int) time.Duration {
	return 1 * time.Second
}

func ExponentialDelayFunc(attempt int) time.Duration {
	exp2 := math.Exp2(float64(attempt))
	return time.Duration(int(exp2))
}
