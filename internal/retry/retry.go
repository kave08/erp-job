package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net"
	"net/url"
	"time"
)

type Policy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
	Jitter         float64
}

type Attempt struct {
	Attempt    int
	StatusCode int
	Error      error
	WillRetry  bool
	Duration   time.Duration
}

type Observer func(Attempt)

func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:    4,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     8 * time.Second,
		Multiplier:     2,
		Jitter:         0.2,
	}
}

func Do(ctx context.Context, policy Policy, observer Observer, fn func(context.Context) (int, error)) (Attempt, error) {
	var last Attempt

	for attemptNumber := 1; attemptNumber <= policy.MaxAttempts; attemptNumber++ {
		startedAt := time.Now()
		statusCode, err := fn(ctx)
		last = Attempt{
			Attempt:    attemptNumber,
			StatusCode: statusCode,
			Error:      err,
			Duration:   time.Since(startedAt),
		}

		if err == nil {
			if observer != nil {
				observer(last)
			}
			return last, nil
		}

		last.WillRetry = attemptNumber < policy.MaxAttempts && IsRetryable(statusCode, err)
		if observer != nil {
			observer(last)
		}
		if !last.WillRetry {
			return last, err
		}

		if err := sleep(ctx, backoffForAttempt(policy, attemptNumber)); err != nil {
			return last, err
		}
	}

	return last, last.Error
}

func IsRetryable(statusCode int, err error) bool {
	switch {
	case statusCode == 408 || statusCode == 429:
		return true
	case statusCode >= 500 && statusCode <= 599:
		return true
	case err == nil:
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var urlErr *url.Error
	return errors.As(err, &urlErr)
}

func backoffForAttempt(policy Policy, attempt int) time.Duration {
	backoff := float64(policy.InitialBackoff) * math.Pow(policy.Multiplier, float64(attempt-1))
	if max := float64(policy.MaxBackoff); backoff > max {
		backoff = max
	}

	jitterWindow := backoff * policy.Jitter
	if jitterWindow > 0 {
		backoff += rand.Float64()*(2*jitterWindow) - jitterWindow
	}

	if backoff < 0 {
		backoff = 0
	}

	return time.Duration(backoff)
}

func sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
