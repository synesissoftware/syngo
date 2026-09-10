package sync_test

import (
	. "github.com/synesissoftware/syngo/sync"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_BoolLatch(t *testing.T) {

	t.Run("NewBoolLatch() succeeds", func(t *testing.T) {

		_ = NewBoolLatch()

		assert.True(t, true)
	})

	t.Run("Load() without Set()", func(t *testing.T) {

		latch := NewBoolLatch()

		require.Equal(t, false, latch.Load())
	})

	t.Run("Load() before and after Set()", func(t *testing.T) {

		var flipped bool

		latch := NewBoolLatch()

		require.False(t, latch.Load())
		require.False(t, latch.Load())
		require.False(t, latch.Load())
		require.False(t, latch.Load())

		flipped = latch.Set()
		require.True(t, flipped)

		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())

		flipped = latch.Set()
		require.False(t, flipped)

		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())

		flipped = latch.Set()
		require.False(t, flipped)

		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())
		require.True(t, latch.Load())
	})

	t.Run("hitting Load() from many goroutines, and waiting to Set() until one of the readers hits half the number of loads", func(t *testing.T) {

		latch := NewBoolLatch()

		const numGoroutines = 10
		const numLoads = 2_000_000
		const permitThreshold = numLoads / 2
		const totalLoadCount uint64 = numGoroutines * numLoads

		var wg sync.WaitGroup
		var count atomic.Uint64
		var permitSet atomic.Bool

		falseReads := make([]int64, numGoroutines)
		trueReads := make([]int64, numGoroutines)

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				runtime.Gosched()

				if permitSet.Load() {
					flipped := latch.Set()

					require.True(t, flipped)

					return
				}
			}
		}()

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				for j := 0; j != numLoads; j++ {

					if latch.Load() {

						atomic.AddInt64(&trueReads[i], 1)
					} else {

						atomic.AddInt64(&falseReads[i], 1)
					}

					count.Add(1)

					if j == permitThreshold {
						permitSet.Store(true)
					}
				}
			}(i)
		}

		wg.Wait()

		var numFalseReads int64
		var numTrueReads int64

		require.Equal(t, totalLoadCount, count.Load())

		for i := 0; i != numGoroutines; i++ {
			numFalseReads += falseReads[i]
			numTrueReads += trueReads[i]
		}

		require.Equal(t, int64(totalLoadCount), numFalseReads+numTrueReads)
	})

	t.Run("concurrent Set() flips exactly once", func(t *testing.T) {

		latch := NewBoolLatch()

		const numGoroutines = 64
		const numAttempts = 10_000

		var wg sync.WaitGroup
		var flips atomic.Int64

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != numAttempts; j++ {

					if latch.Set() {
						flips.Add(1)
					}
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())
		require.True(t, latch.Load())
	})

	t.Run("concurrent Set() and Load() settle latched", func(t *testing.T) {

		latch := NewBoolLatch()

		const numSetters = 8
		const numLoaders = 16
		const attempts = 50_000

		var wg sync.WaitGroup
		var flips atomic.Int64
		var sawTrue atomic.Bool

		for i := 0; i != numSetters; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != attempts; j++ {

					if latch.Set() {
						flips.Add(1)
					}
				}
			}()
		}

		for i := 0; i != numLoaders; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != attempts; j++ {

					if latch.Load() {
						sawTrue.Store(true)
					}
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())
		require.True(t, latch.Load())
		require.True(t, sawTrue.Load())
	})
}

func Test_DownLatch(t *testing.T) {

	t.Run("NewDownLatch() succeeds", func(t *testing.T) {

		_ = NewDownLatch(1, 0)

		assert.True(t, true)
	})

	t.Run("Load() without Set()", func(t *testing.T) {

		latch := NewDownLatch(1, 0)

		isLatched, count := latch.Load()

		require.Equal(t, false, isLatched)
		require.Equal(t, int64(1), count)
	})

	t.Run("Load() before and after Set(), for a latch of range [10, 8)", func(t *testing.T) {

		latch := NewDownLatch(10, 8)

		var flipped bool
		var isLatched bool
		var count int64

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(10), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(10), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.False(t, flipped)
			require.False(t, isLatched)
			require.Equal(t, int64(9), count)
		}

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(9), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(9), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.True(t, flipped)
			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}

		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}
		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.False(t, flipped)
			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}

		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}
		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(8), count)
		}
	})

	t.Run("Load() before and after Set(), for a latch of range [2, -1)", func(t *testing.T) {

		latch := NewDownLatch(2, -1)

		var flipped bool
		var isLatched bool
		var count int64

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(2), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(2), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.False(t, flipped)
			require.False(t, isLatched)
			require.Equal(t, int64(1), count)
		}

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(1), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(1), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.False(t, flipped)
			require.False(t, isLatched)
			require.Equal(t, int64(0), count)
		}

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(0), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(0), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.True(t, flipped)
			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}

		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}
		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}

		{
			flipped, isLatched, count = latch.Step()

			require.False(t, flipped)
			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}

		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}
		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			require.Equal(t, int64(-1), count)
		}
	})

	t.Run("NewDownLatch() panics when initialValue <= threshold", func(t *testing.T) {

		require.Panics(t, func() {
			_ = NewDownLatch(0, 0)
		})
		require.Panics(t, func() {
			_ = NewDownLatch(5, 5)
		})
		require.Panics(t, func() {
			_ = NewDownLatch(1, 2)
		})
	})

	t.Run("concurrent Step() flips exactly once", func(t *testing.T) {

		const distance int64 = 100
		const threshold int64 = 7
		const initialValue = threshold + distance

		latch := NewDownLatch(initialValue, threshold)

		const numGoroutines = 32
		const extraSteps = 1_000

		var wg sync.WaitGroup
		var flips atomic.Int64
		var latchedObservations atomic.Int64

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				// Together the goroutines overshoot the latch distance heavily.
				steps := int(distance)/numGoroutines + extraSteps
				for j := 0; j != steps; j++ {

					flipped, isLatched, count := latch.Step()

					if flipped {
						flips.Add(1)
						require.True(t, isLatched)
						require.Equal(t, threshold, count)
					}

					if isLatched {
						latchedObservations.Add(1)
						require.Equal(t, threshold, count)
					} else {
						require.Greater(t, count, threshold)
					}
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())
		require.Greater(t, latchedObservations.Load(), int64(0))

		isLatched, count := latch.Load()
		require.True(t, isLatched)
		require.Equal(t, threshold, count)
	})

	t.Run("exact-distance concurrent Step() — one flip, rest post-latch", func(t *testing.T) {

		const distance int64 = 200
		const threshold int64 = 0

		latch := NewDownLatch(distance, threshold)

		var wg sync.WaitGroup
		var flips atomic.Int64

		for i := 0; i != int(distance); i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				flipped, isLatched, count := latch.Step()

				if flipped {
					flips.Add(1)
					require.True(t, isLatched)
					require.Equal(t, threshold, count)
				} else if isLatched {
					require.Equal(t, threshold, count)
				} else {
					require.Greater(t, count, threshold)
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())

		isLatched, count := latch.Load()
		require.True(t, isLatched)
		require.Equal(t, threshold, count)
	})

	t.Run("concurrent Load() during Step() stays in range", func(t *testing.T) {

		const initialValue int64 = 5_000
		const threshold int64 = 0

		latch := NewDownLatch(initialValue, threshold)

		var writersWG sync.WaitGroup
		var readersWG sync.WaitGroup
		var writersDone atomic.Bool
		var loadViolation atomic.Bool

		for r := 0; r != 8; r++ {

			readersWG.Add(1)
			go func() {
				defer readersWG.Done()

				for {
					isLatched, count := latch.Load()

					if count > initialValue || count < threshold {
						loadViolation.Store(true)
						return
					}

					if isLatched && count != threshold {
						loadViolation.Store(true)
						return
					}

					if !isLatched && count <= threshold {
						loadViolation.Store(true)
						return
					}

					if writersDone.Load() {
						return
					}

					runtime.Gosched()
				}
			}()
		}

		writersWG.Add(1)
		go func() {
			defer writersWG.Done()

			for i := int64(0); i != initialValue+100; i++ {
				latch.Step()
			}
		}()

		writersWG.Wait()
		writersDone.Store(true)
		readersWG.Wait()

		require.False(t, loadViolation.Load())

		isLatched, count := latch.Load()
		require.True(t, isLatched)
		require.Equal(t, threshold, count)
	})
}

func Test_UpLatch(t *testing.T) {

	t.Run("NewUpLatch() succeeds", func(t *testing.T) {

		_ = NewUpLatch(-1_000_000_000, 1_000_000_000)

		assert.True(t, true)
	})

	t.Run("Load() without Set(), for a latch of range [1, 3)", func(t *testing.T) {

		latch := NewUpLatch(1, 3)

		isLatched, count := latch.Load()

		require.Equal(t, false, isLatched)
		require.Equal(t, int64(1), count)
	})

	t.Run("Load() without Set(), for a latch of range [-101, 102)", func(t *testing.T) {

		latch := NewUpLatch(-101, 102)

		isLatched, count := latch.Load()

		require.Equal(t, false, isLatched)
		require.Equal(t, int64(-101), count)
	})

	t.Run("Load() before and after Set(), for a latch of range [-10, 8)", func(t *testing.T) {

		latch := NewUpLatch(-10, 8)

		var flipped bool
		var isLatched bool
		var count int64

		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			require.Equal(t, int64(-10), count)
		}
		{
			isLatched, count = latch.Load()

			require.False(t, isLatched)
			assert.Equal(t, int64(-10), count)
		}

		for i := 0; i != 17; i++ {

			expectedNewCount := int64(-10 + (1 + i))

			{
				flipped, isLatched, count = latch.Step()

				require.False(t, flipped)
				require.False(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}

			{
				isLatched, count = latch.Load()

				require.False(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}
			{
				isLatched, count = latch.Load()

				require.False(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}
		}

		{
			flipped, isLatched, count = latch.Step()

			require.True(t, flipped)
			require.True(t, isLatched)
			assert.Equal(t, int64(8), count)
		}

		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			assert.Equal(t, int64(8), count)
		}
		{
			isLatched, count = latch.Load()

			require.True(t, isLatched)
			assert.Equal(t, int64(8), count)
		}

		for i := 0; i != 1_000; i++ {

			expectedNewCount := int64(8)

			{
				flipped, isLatched, count = latch.Step()

				require.False(t, flipped)
				require.True(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}

			{
				isLatched, count = latch.Load()

				require.True(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}
			{
				isLatched, count = latch.Load()

				require.True(t, isLatched)
				assert.Equal(t, expectedNewCount, count)
			}
		}
	})

	t.Run("NewUpLatch() panics when initialValue >= threshold", func(t *testing.T) {

		require.Panics(t, func() {
			_ = NewUpLatch(0, 0)
		})
		require.Panics(t, func() {
			_ = NewUpLatch(5, 5)
		})
		require.Panics(t, func() {
			_ = NewUpLatch(3, 1)
		})
	})

	t.Run("concurrent Step() flips exactly once", func(t *testing.T) {

		const initialValue int64 = -20
		const threshold int64 = 80
		const distance = threshold - initialValue

		latch := NewUpLatch(initialValue, threshold)

		const numGoroutines = 32
		const extraSteps = 1_000

		var wg sync.WaitGroup
		var flips atomic.Int64
		var latchedObservations atomic.Int64

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				steps := int(distance)/numGoroutines + extraSteps
				for j := 0; j != steps; j++ {

					flipped, isLatched, count := latch.Step()

					if flipped {
						flips.Add(1)
						require.True(t, isLatched)
						require.Equal(t, threshold, count)
					}

					if isLatched {
						latchedObservations.Add(1)
						require.Equal(t, threshold, count)
					} else {
						require.Less(t, count, threshold)
					}
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())
		require.Greater(t, latchedObservations.Load(), int64(0))

		isLatched, count := latch.Load()
		require.True(t, isLatched)
		require.Equal(t, threshold, count)
	})

	t.Run("exact-distance concurrent Step() — one flip", func(t *testing.T) {

		const initialValue int64 = 0
		const threshold int64 = 250

		latch := NewUpLatch(initialValue, threshold)

		var wg sync.WaitGroup
		var flips atomic.Int64

		for i := 0; i != int(threshold-initialValue); i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				flipped, isLatched, count := latch.Step()

				if flipped {
					flips.Add(1)
					require.True(t, isLatched)
					require.Equal(t, threshold, count)
				} else if isLatched {
					require.Equal(t, threshold, count)
				} else {
					require.Less(t, count, threshold)
					require.Greater(t, count, initialValue)
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(1), flips.Load())

		isLatched, count := latch.Load()
		require.True(t, isLatched)
		require.Equal(t, threshold, count)
	})
}
