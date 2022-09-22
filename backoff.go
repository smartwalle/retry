package retry

import "time"

type Backoff interface {
	Backoff(retries int) time.Duration
}
