//go:build synthetic

package transcript

import (
	"context"
	"testing"
	"time"
)

func TestSyntheticStreamDelayAndOneDisconnect(t *testing.T) {
	t.Setenv("VN_STREAM_DELAY_MS", "100")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if BeforeYield(ctx) != context.Canceled {
		t.Fatal("delay did not honor cancellation")
	}
	t.Setenv("VN_STREAM_DROP_ONCE_AFTER_MS", "1")
	t.Setenv("VN_STREAM_FAIL_AFTER_MS", "")
	streamStarted = time.Now()
	dropScheduled = false
	ch, stop := StreamFailure()
	defer stop()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("disconnect deadline did not fire")
	}
	ch, stop = StreamFailure()
	defer stop()
	if ch != nil {
		t.Fatal("disconnect repeated")
	}
	t.Setenv("VN_STREAM_FAIL_AFTER_MS", "1")
	ch, stop = StreamFailure()
	defer stop()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("permanent failure not immediate after deadline")
	}
}

func TestSyntheticCommandDelayHonorsCancellation(t *testing.T) {
	t.Setenv("VN_COMMAND_DELAY_MS", "250")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if BeforeAccept(ctx) != context.Canceled {
		t.Fatal("command delay ignored request cancellation")
	}
}
