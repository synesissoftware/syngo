package sync_test

import (
	. "github.com/synesissoftware/syngo/sync"

	"github.com/stretchr/testify/require"

	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_HighWaterMark_DEFAULT(t *testing.T) {

	var hwm HighWaterMark

	current, maximum := hwm.Load()

	require.Equal(t, int64(0), current)
	require.Equal(t, int64(0), maximum)
}

func Test_HighWaterMark_Set_and_Load_sequential(t *testing.T) {

	t.Run("single ascending Set", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(123)

		current, maximum := hwm.Load()

		require.Equal(t, int64(123), current)
		require.Equal(t, int64(123), maximum)
	})

	t.Run("ascending sequence raises current and maximum together", func(t *testing.T) {

		var hwm HighWaterMark

		for _, v := range []int64{1, 2, 5, 10, 100, 1_000} {

			hwm.Set(v)

			current, maximum := hwm.Load()

			require.Equal(t, v, current)
			require.Equal(t, v, maximum)
		}
	})

	t.Run("descending after peak retains maximum", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(100)
		hwm.Set(50)
		hwm.Set(25)
		hwm.Set(0)

		current, maximum := hwm.Load()

		require.Equal(t, int64(0), current)
		require.Equal(t, int64(100), maximum)
	})

	t.Run("equal Set does not change maximum", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(42)
		hwm.Set(42)
		hwm.Set(42)

		current, maximum := hwm.Load()

		require.Equal(t, int64(42), current)
		require.Equal(t, int64(42), maximum)
	})

	t.Run("new peak after descent raises maximum again", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(10)
		hwm.Set(1)
		hwm.Set(100)
		hwm.Set(50)

		current, maximum := hwm.Load()

		require.Equal(t, int64(50), current)
		require.Equal(t, int64(100), maximum)
	})

	t.Run("negative values below the zero default leave maximum at zero", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(-1)
		hwm.Set(-100)
		hwm.Set(-50)

		current, maximum := hwm.Load()

		require.Equal(t, int64(-50), current)
		require.Equal(t, int64(0), maximum)
	})

	t.Run("positive peak then negatives retain the positive maximum", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(7)
		hwm.Set(-3)
		hwm.Set(-30)

		current, maximum := hwm.Load()

		require.Equal(t, int64(-30), current)
		require.Equal(t, int64(7), maximum)
	})

	t.Run("extremes of int64", func(t *testing.T) {

		var hwm HighWaterMark

		hwm.Set(math.MinInt64)

		current, maximum := hwm.Load()

		require.Equal(t, int64(math.MinInt64), current)
		require.Equal(t, int64(0), maximum)

		hwm.Set(math.MaxInt64)

		current, maximum = hwm.Load()

		require.Equal(t, int64(math.MaxInt64), current)
		require.Equal(t, int64(math.MaxInt64), maximum)

		hwm.Set(0)

		current, maximum = hwm.Load()

		require.Equal(t, int64(0), current)
		require.Equal(t, int64(math.MaxInt64), maximum)
	})
}

func Test_HighWaterMark_concurrent_Sets_establish_global_maximum(t *testing.T) {

	const numGoroutines = 32
	const numSetsPerGoroutine = 5_000
	const totalSets = numGoroutines * numSetsPerGoroutine

	values := make([]int64, totalSets)
	var expectedMaximum int64

	for i := 0; i != totalSets; i++ {

		// Deterministic, non-monotonic mix so every goroutine races peaks.
		v := int64((i*7_919 + 13) % 100_003)
		if i%17 == 0 {
			v = -v
		}

		values[i] = v

		if v > expectedMaximum {
			expectedMaximum = v
		}
	}

	var hwm HighWaterMark
	var wg sync.WaitGroup

	for g := 0; g != numGoroutines; g++ {

		wg.Add(1)
		go func(g int) {
			defer wg.Done()

			base := g * numSetsPerGoroutine
			for j := 0; j != numSetsPerGoroutine; j++ {
				hwm.Set(values[base+j])
			}
		}(g)
	}

	wg.Wait()

	_, maximum := hwm.Load()

	require.Equal(t, expectedMaximum, maximum)
}

func Test_HighWaterMark_concurrent_Sets_with_shared_peak(t *testing.T) {

	const numGoroutines = 64
	const peak int64 = 1_000_000

	var hwm HighWaterMark
	var wg sync.WaitGroup

	for g := 0; g != numGoroutines; g++ {

		wg.Add(1)
		go func(g int) {
			defer wg.Done()

			// Every goroutine writes the same peak amid lower values.
			hwm.Set(int64(g))
			hwm.Set(peak)
			hwm.Set(int64(-g))
		}(g)
	}

	wg.Wait()

	_, maximum := hwm.Load()

	require.Equal(t, peak, maximum)
}

func Test_HighWaterMark_concurrent_Load_observes_monotonic_maximum(t *testing.T) {

	const numWriters = 16
	const numReaders = 16
	const numSetsPerWriter = 10_000

	var hwm HighWaterMark
	var writersWG sync.WaitGroup
	var readersWG sync.WaitGroup
	var writersDone atomic.Bool
	var monotonicViolation atomic.Bool
	var expectedMaximum atomic.Int64

	for w := 0; w != numWriters; w++ {

		writersWG.Add(1)
		go func(w int) {
			defer writersWG.Done()

			base := int64(w * numSetsPerWriter)
			for j := 0; j != numSetsPerWriter; j++ {

				v := base + int64(j)
				hwm.Set(v)

				for {
					prev := expectedMaximum.Load()
					if v <= prev {
						break
					}
					if expectedMaximum.CompareAndSwap(prev, v) {
						break
					}
				}
			}
		}(w)
	}

	for r := 0; r != numReaders; r++ {

		readersWG.Add(1)
		go func() {
			defer readersWG.Done()

			var lastMaximum int64

			for {
				_, maximum := hwm.Load()

				if maximum < lastMaximum {
					monotonicViolation.Store(true)
					return
				}

				lastMaximum = maximum

				if writersDone.Load() {
					_, maximum = hwm.Load()
					if maximum < lastMaximum {
						monotonicViolation.Store(true)
					}
					return
				}

				runtime.Gosched()
			}
		}()
	}

	writersWG.Wait()
	writersDone.Store(true)
	readersWG.Wait()

	require.False(t, monotonicViolation.Load())

	_, maximum := hwm.Load()
	require.Equal(t, expectedMaximum.Load(), maximum)
}

func Test_HighWaterMark_concurrent_mixed_ascend_descend(t *testing.T) {

	const numGoroutines = 24
	const rounds = 2_000
	const peak int64 = 50_000

	var hwm HighWaterMark
	var wg sync.WaitGroup

	for g := 0; g != numGoroutines; g++ {

		wg.Add(1)
		go func(g int) {
			defer wg.Done()

			for r := 0; r != rounds; r++ {

				hwm.Set(int64(r % 100))
				hwm.Set(peak - int64(g))
				hwm.Set(peak)
				hwm.Set(int64(-(r + g)))
			}
		}(g)
	}

	wg.Wait()

	current, maximum := hwm.Load()

	require.Equal(t, peak, maximum)
	// current is whatever the last Set stored; must be one of the written
	// forms, and cannot exceed the recorded maximum once writers have quit.
	require.LessOrEqual(t, current, maximum)
}

func Test_HighWaterMark_final_Set_after_contention_wins_current(t *testing.T) {

	const numGoroutines = 32
	const numSetsPerGoroutine = 2_000

	var hwm HighWaterMark
	var wg sync.WaitGroup

	for g := 0; g != numGoroutines; g++ {

		wg.Add(1)
		go func(g int) {
			defer wg.Done()

			for j := 0; j != numSetsPerGoroutine; j++ {
				hwm.Set(int64(g*numSetsPerGoroutine + j))
			}
		}(g)
	}

	wg.Wait()

	const finalCurrent int64 = -999
	hwm.Set(finalCurrent)

	current, maximum := hwm.Load()

	require.Equal(t, finalCurrent, current)
	require.Equal(t, int64(numGoroutines*numSetsPerGoroutine-1), maximum)
}
