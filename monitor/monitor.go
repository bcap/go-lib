package monitor

import (
	"context"
	"time"
)

// Monitor runs a function at a given interval until the context is done.
// The function is called with the time passed since the start of the monitor
// and a boolean indicating if the monitor is done and thus this is the last call
func Monitor(ctx context.Context, interval time.Duration, fn func(timePassed time.Duration, end bool)) {
	start := time.Now()
	go func() {
		run := true
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for run {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				run = false
			}
			fn(time.Since(start), !run)
		}
	}()
}
