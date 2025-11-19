package retry_test

import (
	"context"
	"errors"
	"github.com/smartwalle/retry"
	"testing"
	"time"
)

type ConstantBackoff struct {
	RetryDelay int
}

func (b *ConstantBackoff) Backoff(retries int) time.Duration {
	return time.Duration(b.RetryDelay) * time.Second
}

func (b *ConstantBackoff) ShouldRetry(err error) bool {
	return true
}

func TestRetry(t *testing.T) {
	var backoff = &ConstantBackoff{
		RetryDelay: 1,
	}

	var tests = []struct {
		MaxAttempts int
	}{
		{MaxAttempts: -10},
		{MaxAttempts: -1},
		{MaxAttempts: 0},
		{MaxAttempts: 1},
		{MaxAttempts: 2},
		{MaxAttempts: 100},
	}

	for _, test := range tests {
		var actual, err = retry.Do(context.Background(), backoff, test.MaxAttempts, func(ctx context.Context, attempt int) (int, error) {
			return test.MaxAttempts, nil
		})
		if err != nil {
			t.Fatalf("返回错误: %+v", err)
		}
		if actual != test.MaxAttempts {
			t.Fatalf("期望值为: %d，实际值为: %d", test.MaxAttempts, actual)
		}
	}
}

func TestRetryTimeout(t *testing.T) {
	var backoff = &ConstantBackoff{
		RetryDelay: 1,
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	var _, err = retry.Do(ctx, backoff, 3, func(ctx context.Context, attempt int) (int, error) {
		return 0, errors.New("failed")
	})

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("应该返回超时错误，实际错误为: %+v", err)
	}
}

func TestRetryMaxAttempts(t *testing.T) {
	var backoff = &ConstantBackoff{
		RetryDelay: 1,
	}

	var tests = []struct {
		MaxAttempts int
	}{
		{MaxAttempts: 0},
		{MaxAttempts: 1},
		{MaxAttempts: 2},
		{MaxAttempts: 3},
	}

	for _, test := range tests {
		var actual, err = retry.Do(context.Background(), backoff, test.MaxAttempts, func(ctx context.Context, attempt int) (int, error) {
			if attempt < test.MaxAttempts {
				return 0, errors.New("failed")
			}
			return test.MaxAttempts, nil
		})
		if err != nil {
			t.Fatalf("返回错误: %+v", err)
		}
		if actual != test.MaxAttempts {
			t.Fatalf("期望值为: %d，实际值为: %d", test.MaxAttempts, actual)
		}
	}
}
