package retry

import (
	"context"
	"time"
)

type Strategy interface {
	Backoff(ctx context.Context, attempt int) time.Duration

	ShouldRetry(ctx context.Context, err error) bool
}
