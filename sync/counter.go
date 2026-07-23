// Copyright 2019-2025 Harold Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 15th November 2025
 * Updated:
 */

// Definition of a number of types that have one-way behaviour.

package sync

import (
	"sync/atomic"
)

// Count-down-to-zero (and beyond, a bit) numeric latch
type _baseCounter struct {
	value int64
}

func (l *_baseCounter) step(increment int64) (newCount int64) {

	newCount = atomic.AddInt64(&l.value, increment)

	return
}

func (l *_baseCounter) load() (count int64) {

	count = atomic.LoadInt64(&l.value)

	return
}

// A unidirectional latch that counts down from an initial value to a lower
// threshold that may be operated safely by multiple concurrent goroutines.
type DownCounter struct {
	_baseCounter
}

// Creates a new DownLatch.
//
// Preconditions:
// - initialValue > threshold;
// - initialValue - threshold <= MaxLatchDistance;
func NewDownCounter(initialValue int64) DownCounter {

	return DownCounter{
		_baseCounter: _baseCounter{
			value: initialValue,
		},
	}
}

func (l *DownCounter) Step() (newCount int64) {

	newCount = l._baseCounter.step(-1)

	return
}

// Obtains the current value of the latch, without changing its state.
func (l *DownCounter) Load() (count int64) {

	count = l._baseCounter.load()

	return
}

// A unidirectional latch that counts up from an initial value to a higher
// threshold that may be operated safely by multiple concurrent goroutines.
type UpCounter struct {
	_baseCounter
}

// Creates a new UpLatch.
//
// Preconditions:
// - initialValue < threshold;
// - threshold - initialValue <= MaxLatchDistance;
func NewUpCounter(initialValue int64) UpCounter {

	return UpCounter{
		_baseCounter: _baseCounter{
			value: initialValue,
		},
	}
}

func (l *UpCounter) Step() (count int64) {

	count = l._baseCounter.step(1)

	return
}

// Obtains the current value of the latch, without changing its state.
func (l *UpCounter) Load() (count int64) {

	count = l._baseCounter.load()

	return
}
