//go:build synthetic

package transcript

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"
)

var streamStarted = time.Now()
var dropMu sync.Mutex
var dropScheduled bool

func milliseconds(name string) time.Duration {
	n, _ := strconv.Atoi(os.Getenv(name))
	if n <= 0 {
		return 0
	}
	return time.Duration(n) * time.Millisecond
}
func BeforeYield(ctx context.Context) error {
	delay := milliseconds("VN_STREAM_DELAY_MS")
	if delay == 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func StreamFailure() (<-chan time.Time, func()) {
	deadline := milliseconds("VN_STREAM_FAIL_AFTER_MS")
	if deadline == 0 {
		deadline = milliseconds("VN_STREAM_DROP_ONCE_AFTER_MS")
		if deadline == 0 {
			return nil, func() {}
		}
		dropMu.Lock()
		if dropScheduled {
			dropMu.Unlock()
			return nil, func() {}
		}
		dropScheduled = true
		dropMu.Unlock()
	}
	remaining := time.Until(streamStarted.Add(deadline))
	if remaining < 0 {
		remaining = 0
	}
	timer := time.NewTimer(remaining)
	return timer.C, func() { timer.Stop() }
}

func BeforeAccept(ctx context.Context) error {
	delay := milliseconds("VN_COMMAND_DELAY_MS")
	if delay == 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
