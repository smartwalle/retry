package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, backoff Backoff, max int, fn func(context.Context, int) (T, error)) (T, error) {
	var err error
	var value T
	for i := 0; i <= max; i++ {
		if i > 0 {
			var delay = backoff.Duration(i)
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

		var ctxErr = ctx.Err()
		if ctxErr != nil {
			return value, ctxErr
		}
	}
	return value, err
}
