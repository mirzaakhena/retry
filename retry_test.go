package retry

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type stringErrorType string

func (s stringErrorType) Error() string {
	return "stringErrorType"
}

type structErrorType struct{}

func (s structErrorType) Error() string {
	return "structErrorType"
}

type compositionErrorType struct{ error }

// errorProducer is an object that can produce an error several times when called
type errorProducer struct {
	MessageWhenCalled    string
	ErrorToProduce       error
	CounterBeforeSuccess int
	CurrentCounter       int
}

func newErrorProducer(counterBeforeSuccess int, errorToProduce error, message string) *errorProducer {
	return &errorProducer{
		MessageWhenCalled:    message,
		ErrorToProduce:       errorToProduce,
		CounterBeforeSuccess: counterBeforeSuccess,
		CurrentCounter:       0,
	}
}

// Run will produce an error depend on setting
func (r *errorProducer) Run() error {
	r.CurrentCounter++
	fmt.Printf("%s %d : ", r.MessageWhenCalled, r.CurrentCounter)
	if r.CurrentCounter <= r.CounterBeforeSuccess {
		fmt.Println("fail")
		return r.ErrorToProduce
	}
	fmt.Println("SUCCESS")
	return nil
}

// TestErrorProducer_001 is the ErrorProducer it self
func TestErrorProducer_001(t *testing.T) {

	{
		e := newErrorProducer(3, fmt.Errorf("direct error fmt"), "test")
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Nil(t, e.Run())
	}

	{
		e := newErrorProducer(5, stringErrorType("2"), "test")
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Error(t, e.Run())
		assert.Nil(t, e.Run())
	}

	{
		e := newErrorProducer(0, structErrorType{}, "test")
		assert.Nil(t, e.Run())
	}

	{
		e := newErrorProducer(1, compositionErrorType{}, "test")
		assert.Error(t, e.Run())
		assert.Nil(t, e.Run())
	}

}

// TestRetry_002 should run 3 times because of max retry and listen to any error
func TestRetry_002(t *testing.T) {

	// error will occur 5 times, success is in the 6th attempts
	e := newErrorProducer(5, fmt.Errorf("direct error fmt"), "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 3,
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: DefaultAnyError(), // <--- any error are welcome
	})
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	assert.Nil(t, err)
	r.OnRetry(func() error { return e.Run() })
}

// TestRetry_003 retry should run once since the error type does not match
func TestRetry_003(t *testing.T) {
	e := newErrorProducer(3, fmt.Errorf("direct error fmt"), "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 3,
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: []error{stringErrorType("")}, // <--- only this type of error
	})
	assert.Nil(t, err)
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	r.OnRetry(func() error { return e.Run() })
}

// TestRetry_0031 retry should run once since the error type does not match
func TestRetry_0031(t *testing.T) {
	e := newErrorProducer(3, fmt.Errorf("direct error fmt"), "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 3,
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: []error{structErrorType{}}, // <--- only this type of error
	})
	assert.Nil(t, err)
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	r.OnRetry(func() error { return e.Run() })
}

// TestRetry_004 retry should only error 3 times and in the fourth times, it will success (no need to retry)
func TestRetry_004(t *testing.T) {
	e := newErrorProducer(3, fmt.Errorf("direct error fmt"), "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 10, // <--- we set Max Retry to 10
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: DefaultAnyError(),
	})
	assert.Nil(t, err)
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	r.OnRetry(func() error { return e.Run() })
}

// TestRetry_005 retry should only error 3 times and in the fourth times, it will success (no need to retry)
func TestRetry_005(t *testing.T) {
	e := newErrorProducer(3, structErrorType{}, "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 10, // <--- we set Max Retry to 10
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: []error{structErrorType{}},
	})
	assert.Nil(t, err)
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	r.OnRetry(func() error { return e.Run() })
}

// TestRetry_006 retry should only error 3 times and in the fourth times, it will success (no need to retry)
func TestRetry_006(t *testing.T) {
	e := newErrorProducer(3, compositionErrorType{}, "test")

	r := New()
	err := r.SetConfig(Config{
		MaxRetry: 10, // <--- we set Max Retry to 10
		DelayFunc: func(attempt int) time.Duration {
			return 1 * time.Second
		},
		ByErrors: []error{stringErrorType(""), compositionErrorType{}},
	})
	assert.Nil(t, err)
	r.SetRetryListener(func(info Info) {
		fmt.Printf("%v %v\n", info.RetryCount, info.Error)
	})
	r.OnRetry(func() error { return e.Run() })
}
