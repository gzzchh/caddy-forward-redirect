package followredirect

import (
	"fmt"
	"sync/atomic"
)

type upstreamHost struct {
	numRequests int64 // must be 64-bit aligned on 32-bit systems (see https://golang.org/pkg/sync/atomic/#pkg-note-BUG)
	fails       int64
	unhealthy   int32
}

// NumRequests returns the number of active requests to the upstream.
func (uh *upstreamHost) NumRequests() int {
	return int(atomic.LoadInt64(&uh.numRequests))
}

// Fails returns the number of recent failures with the upstream.
func (uh *upstreamHost) Fails() int {
	return int(atomic.LoadInt64(&uh.fails))
}

// Unhealthy returns whether the upstream is healthy.
func (uh *upstreamHost) Unhealthy() bool {
	return atomic.LoadInt32(&uh.unhealthy) == 1
}

// CountRequest mutates the active request count by
// delta. It returns an error if the adjustment fails.
func (uh *upstreamHost) CountRequest(delta int) error {
	result := atomic.AddInt64(&uh.numRequests, int64(delta))
	if result < 0 {
		return fmt.Errorf("count below 0: %d", result)
	}
	return nil
}

// CountFail mutates the recent failures count by
// delta. It returns an error if the adjustment fails.
func (uh *upstreamHost) CountFail(delta int) error {
	result := atomic.AddInt64(&uh.fails, int64(delta))
	if result < 0 {
		return fmt.Errorf("count below 0: %d", result)
	}
	return nil
}

// SetHealthy sets the upstream has healthy or unhealthy
// and returns true if the new value is different.
func (uh *upstreamHost) SetHealthy(healthy bool) (bool, error) {
	var unhealthy, compare int32 = 1, 0
	if healthy {
		unhealthy, compare = 0, 1
	}
	swapped := atomic.CompareAndSwapInt32(&uh.unhealthy, compare, unhealthy)
	return swapped, nil
}
