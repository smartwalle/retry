package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, backoff Backoff, maxAttempts int, fn func(ctx context.Context, attempt int) (T, error)) (value T, err error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		var delay = time.Duration(0)
		if attempt > 0 {
			delay = backoff.Delay(attempt)
		}
		if delay > 0 {
			var timer = time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return value, ctx.Err()
			case <-timer.C:
			}
		}

		if err = ctx.Err(); err != nil {
			return value, err
		}

		if value, err = fn(ctx, attempt); err == nil {
			return value, nil
		}
	}
	return value, err
}
