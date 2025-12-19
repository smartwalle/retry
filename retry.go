package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, strategy Strategy, maxAttempts int, fn func(ctx context.Context, attempt int) (T, error)) (value T, err error) {
	var attempt = 1
	for {
		if err = ctx.Err(); err != nil {
			return value, err
		}

		if value, err = fn(ctx, attempt); err == nil {
			return value, nil
		}

		attempt += 1

		if attempt > maxAttempts {
			return value, err
		}

		if !strategy.ShouldRetry(err) {
			return value, err
		}

		if err = delay(ctx, strategy.Backoff(attempt)); err != nil {
			return value, err
		}
	}
}

func delay(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	var timer = time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
