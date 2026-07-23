package sync_test

import (
	. "github.com/synesissoftware/syngo/sync"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sync"
	"testing"
)

func Test_DownCounter(t *testing.T) {

	t.Run("NewDownCounter() succeeds", func(t *testing.T) {

		_ = NewDownCounter(1)

		assert.True(t, true)
	})

	t.Run("Load() without Set()", func(t *testing.T) {

		counter := NewDownCounter(1)

		count := counter.Load()

		require.Equal(t, int64(1), count)
	})

	t.Run("Load() before and after step()", func(t *testing.T) {

		counter := NewDownCounter(2)

		var count int64

		{
			count = counter.Load()

			require.Equal(t, int64(2), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(2), count)
		}
		{
			count = counter.Step()

			require.Equal(t, int64(1), count)
		}

		{
			count = counter.Load()

			require.Equal(t, int64(1), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(1), count)
		}

		{
			count = counter.Step()

			require.Equal(t, int64(0), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(0), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(0), count)
		}

		{
			count = counter.Step()

			require.Equal(t, int64(-1), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(-1), count)
		}
		{
			count = counter.Load()

			require.Equal(t, int64(-1), count)
		}
	})

	t.Run("hitting Load() and Step() from many goroutines", func(t *testing.T) {

		counter := NewDownCounter(0)

		const numGoroutines = 10
		const numLoads = 10_000
		const totalLoadCount int64 = numGoroutines * numLoads

		var wg sync.WaitGroup

		for i := 0; i != numGoroutines; i++ {

			wg.Go(func() {

				for j := 0; j != numLoads; j++ {

					_ = counter.Step()
				}
			})
		}

		wg.Wait()

		require.Equal(t, -totalLoadCount, counter.Load())
	})
}

func Test_UpCounter(t *testing.T) {

	t.Run("NewUpCounter() succeeds", func(t *testing.T) {

		_ = NewUpCounter(1)

		assert.True(t, true)
	})

	t.Run("Load() without Set()", func(t *testing.T) {

		counter := NewUpCounter(1)

		count := counter.Load()

		require.Equal(t, int64(1), count)
	})

	t.Run("Load() before and after Set()", func(t *testing.T) {

		counter := NewUpCounter(-10)

		var count int64

		{
			count = counter.Load()

			require.Equal(t, int64(-10), count)
		}
		{
			count = counter.Load()

			assert.Equal(t, int64(-10), count)
		}

		for i := 0; i != 17; i++ {

			expectedNewCount := int64(-10 + (1 + i))

			{
				count = counter.Step()

				assert.Equal(t, expectedNewCount, count)
			}

			{
				count = counter.Load()

				assert.Equal(t, expectedNewCount, count)
			}
			{
				count = counter.Load()

				assert.Equal(t, expectedNewCount, count)
			}
		}

		{
			count = counter.Step()

			assert.Equal(t, int64(8), count)
		}

		{
			count = counter.Load()

			assert.Equal(t, int64(8), count)
		}
		{
			count = counter.Load()

			assert.Equal(t, int64(8), count)
		}
	})

	t.Run("hitting Load() and Step() from many goroutines", func(t *testing.T) {

		counter := NewUpCounter(0)

		const numGoroutines = 10
		const numLoads = 10_000
		const totalLoadCount int64 = numGoroutines * numLoads

		var wg sync.WaitGroup

		for i := 0; i != numGoroutines; i++ {

			wg.Go(func() {

				for j := 0; j != numLoads; j++ {

					_ = counter.Step()
				}
			})
		}

		wg.Wait()

		require.Equal(t, totalLoadCount, counter.Load())
	})
}
