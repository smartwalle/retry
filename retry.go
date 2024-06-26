package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, strategy BackoffStrategy, maxRetries int, fn func(context.Context, int) (T, error)) (T, error) {
	var err error
	var value T
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			var delay = strategy.Duration(i)
			var timer = time.NewTimer(delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				return value, ctx.Err()
			case <-timer.C:
			}
		}

		value, err = fn(ctx, i)
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
