package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, strategy BackoffStrategy, maxAttempts int, fn func(ctx context.Context, attempt int) (T, error)) (T, error) {
	var err error
	var value T
	for attempt := 0; attempt <= maxAttempts; attempt++ {
		if attempt > 0 {
			var delay = strategy.Duration(attempt)
			var timer = time.NewTimer(delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				return value, ctx.Err()
			case <-timer.C:
			}
		}

		value, err = fn(ctx, attempt)
		if err == nil {
			return value, nil
		}

		var nErr = ctx.Err()
		if nErr != nil {
			return value, nErr
		}
	}
	return value, err
}
