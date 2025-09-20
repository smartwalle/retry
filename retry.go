package retry

import (
	"context"
	"time"
)

func Do(ctx context.Context, backoff Backoff, maxAttempts int, fn func(ctx context.Context, attempt int) error) (err error) {
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
				return ctx.Err()
			case <-timer.C:
			}
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		if err = fn(ctx, attempt); err == nil {
			return nil
		}
	}
	return err
}
