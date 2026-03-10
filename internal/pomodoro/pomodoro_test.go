package pomodoro

import (
	"context"
	"testing"
	"time"
)

func TestWaitWithCountdown_completes(t *testing.T) {
	ctx := context.Background()
	result := waitWithCountdown(ctx, 50*time.Millisecond, "T")
	if !result {
		t.Error("expected true when duration elapses, got false")
	}
}

func TestWaitWithCountdown_cancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := waitWithCountdown(ctx, 10*time.Second, "T")
	if result {
		t.Error("expected false when context is cancelled, got true")
	}
}

func TestWaitWithCountdown_cancelMidway(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	result := waitWithCountdown(ctx, 10*time.Second, "T")
	if result {
		t.Error("expected false when context cancelled mid-countdown, got true")
	}
}
