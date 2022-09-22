package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, backoff Backoff, max int, fn func(context.Context) (T, error)) (T, error) {
	var err error
	var value T
	for i := 0; i <= max; i++ {
		if i > 0 {
			var delay = backoff.Backoff(i)
			var timer = time.NewTimer(delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				return value, ctx.Err()
			case <-timer.C:
			}
		}

		value, err = fn(ctx)
		if err == nil {
			return value, nil
		}

		if ctx.Err() != nil {
			return value, ctx.Err()
		}
	}
	return value, err
}
