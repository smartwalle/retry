package retry

import "time"

type Backoff interface {
	Delay(retries int) time.Duration

	ShouldRetry(err error) bool
}
