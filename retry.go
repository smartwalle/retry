package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, backoff Backoff, maxAttempts int, fn func(ctx context.Context, attempt int) (T, error)) (value T, err error) {
	var attempt = 0
	for {
		if value, err = fn(ctx, attempt); err == nil {
			return value, nil
		}

		attempt += 1

		if attempt > maxAttempts {
			return value, err
		}

		if !backoff.ShouldRetry(err) {
			return value, err
		}

		var delay = backoff.Delay(attempt)
		if delay > 0 {
			var timer = time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return value, ctx.Err()
			}
		}

		if err = ctx.Err(); err != nil {
			return value, err
		}
	}
}
