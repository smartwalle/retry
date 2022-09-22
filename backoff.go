package retry

import "time"

type Backoff interface {
	Duration(retries int) time.Duration
}
