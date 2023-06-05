package retry

import "time"

type BackoffStrategy interface {
	Duration(retries int) time.Duration
}
