package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/smartwalle/retry"
	"time"
)

func main() {
	var ctx, _ = context.WithTimeout(context.Background(), time.Second*10)

	var value, err = retry.Do[int](ctx, &Backoff{}, 3, func(ctx context.Context, retries int) (int, error) {
		fmt.Println(retries, time.Now())
		return 10, errors.New("sss")
	})

	fmt.Println(value, err)
}

type Backoff struct {
}

func (*Backoff) Duration(retries int) time.Duration {
	return time.Second * time.Duration(retries)
}
