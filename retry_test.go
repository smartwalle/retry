package retry_test

import (
	"context"
	"errors"
	"github.com/smartwalle/retry"
	"testing"
	"time"
)

type ConstantBackoff struct {
	Delay int
}

func (b *ConstantBackoff) Backoff(attempt int) time.Duration {
	return time.Duration(b.Delay) * time.Second
}

func (b *ConstantBackoff) ShouldRetry(err error) bool {
	return true
}

func TestRetryTimeout(t *testing.T) {
	var backoff = &ConstantBackoff{
		Delay: 3,
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var _, err = retry.Do(ctx, backoff, 3, func(ctx context.Context, attempt int) (int, error) {
		return 0, errors.New("failed")
	})

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("应该返回超时错误，实际错误为: %+v", err)
	}
}
