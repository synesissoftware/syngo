package sync

import (
	"sync/atomic"
)

// Thread-safe level that remembers the maximum value assigned to it.
type HighWaterMark struct {
	current int64
	maximum int64
}

// Obtains the current and maximum values.
//
// # Note:
// It is not guaranteed that current may not be larger than maximum.
func (hwm *HighWaterMark) Load() (
	current int64,
	maximum int64,
) {

	current = atomic.LoadInt64(&hwm.current)
	maximum = atomic.LoadInt64(&hwm.maximum)

	return
}

// Sets the current value and, if larger than the current maximum, updates
// that accordingly.
//
// # Note:
// A given thread may observe a maximum value that is larger than the value
// it just set, in the case that another thread has already set a larger
// value.
func (hwm *HighWaterMark) Set(value int64) {

	previousCurrent := atomic.SwapInt64(&hwm.current, value)

	if previousCurrent >= value {

		// nothing to do, since when the previous `hwm.current` was set it would
		// have checked the maximum
	} else {

		for {

			maximum := atomic.LoadInt64(&hwm.maximum)

			if maximum >= value {

				// nothing to do, since this `maximum` lower than a prior
				// `hwm.maximum`

				break
			} else {

				if atomic.CompareAndSwapInt64(&hwm.maximum, maximum, value) {

					// `hwm.maximum` was not changed since we read it, so we replace
					// with `value`

					break
				} else {

					// failed to replace, because
				}
			}
		}
	}
}
