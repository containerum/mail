package utils

import (
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(*StopRetry); ok {
			// Return the original error for later checking
			return s.Err
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type StopRetry struct {
	Err error
}

func (e *StopRetry) Error() string {
	return e.Err.Error()
}
