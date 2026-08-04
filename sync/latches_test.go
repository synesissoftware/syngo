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
}
