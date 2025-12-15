package retry

import "time"

type Strategy interface {
	Backoff(attempt int) time.Duration

	ShouldRetry(err error) bool
}
