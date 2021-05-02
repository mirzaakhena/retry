package retry

import (
	"errors"
	"time"
)

type Info struct {
	Error      interface{}
	RetryCount int
}

type Retry struct {
	config        Config
	retryCount    int
	retryListener func(info Info)
}

func New() *Retry {
	return &Retry{
		config: DefaultConfig(),
	}
}

func (r *Retry) SetConfig(c Config) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	r.config = c
	return nil
}

func (r *Retry) OnRetry(retryableFunc func() error) {

	r.retryCount = 0

	for r.retryCount < r.config.MaxRetry {

		// call the func that need to be retried
		err := retryableFunc()

		// do panic handler first
		if pErr := recover(); pErr != nil {

			// call listener if registered
			r.callListener(pErr)

			// TODO do you have any suggestion what should we do?
			// currently we just stop from retry
			return
		}

		// no error just exit
		if err == nil {
			return
		}

		// if no match any error listed then no need to retry
		if !r.isMatchListedErrors(err) {

			// call listener if registered
			r.callListener(err)

			return
		}

		// for the last try, no need to wait anymore, just exit
		if !r.isTheLastRetry() {
			time.Sleep(r.config.DelayFunc(r.retryCount))
		}

		r.retryCount++
	}

	return
}

func (r *Retry) SetRetryListener(reporter func(info Info)) {
	r.retryListener = reporter
}

func (r *Retry) callListener(pErr interface{}) {
	if r.retryListener == nil {
		return
	}
	r.retryListener(Info{Error: pErr, RetryCount: r.retryCount})
}

func (r *Retry) isTheLastRetry() bool {
	return r.retryCount >= r.config.MaxRetry-1
}

func (r *Retry) isMatchListedErrors(err error) bool {
	if errors.Is(r.config.ByErrors[0], AnyError{}) {
		return true
	}
	for _, testErr := range r.config.ByErrors {
		if errors.Is(err, testErr) {
			return true
		}
	}
	return false
}
