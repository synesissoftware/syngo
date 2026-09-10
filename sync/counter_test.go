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

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != numLoads; j++ {

					_ = counter.Step()
				}
			}()
		}

		wg.Wait()

		require.Equal(t, -totalLoadCount, counter.Load())
	})

	t.Run("concurrent Step() return values are unique", func(t *testing.T) {

		const numGoroutines = 16
		const stepsPerGoroutine = 2_000
		const totalSteps = numGoroutines * stepsPerGoroutine

		counter := NewDownCounter(0)

		results := make([]int64, totalSteps)
		var wg sync.WaitGroup
		var next atomic.Int64

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != stepsPerGoroutine; j++ {

					got := counter.Step()
					slot := next.Add(1) - 1
					results[slot] = got
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(totalSteps), next.Load())
		require.Equal(t, int64(-totalSteps), counter.Load())

		seen := make(map[int64]struct{}, totalSteps)
		for _, v := range results {
			_, exists := seen[v]
			require.False(t, exists, "duplicate Step() return %d", v)
			seen[v] = struct{}{}
		}

		require.Equal(t, totalSteps, len(seen))
	})

	t.Run("concurrent Load() interleaved with Step()", func(t *testing.T) {

		const numSteppers = 8
		const numLoaders = 8
		const stepsPerStepper = 5_000
		const totalSteps = numSteppers * stepsPerStepper

		counter := NewDownCounter(0)

		var steppersWG sync.WaitGroup
		var loadersWG sync.WaitGroup
		var steppersDone atomic.Bool

		for i := 0; i != numSteppers; i++ {

			steppersWG.Add(1)
			go func() {
				defer steppersWG.Done()

				for j := 0; j != stepsPerStepper; j++ {
					_ = counter.Step()
				}
			}()
		}

		for i := 0; i != numLoaders; i++ {

			loadersWG.Add(1)
			go func() {
				defer loadersWG.Done()

				for !steppersDone.Load() {
					_ = counter.Load()
					runtime.Gosched()
				}
			}()
		}

		steppersWG.Wait()
		steppersDone.Store(true)
		loadersWG.Wait()

		require.Equal(t, int64(-totalSteps), counter.Load())
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

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != numLoads; j++ {

					_ = counter.Step()
				}
			}()
		}

		wg.Wait()

		require.Equal(t, totalLoadCount, counter.Load())
	})

	t.Run("concurrent Step() return values are unique", func(t *testing.T) {

		const numGoroutines = 16
		const stepsPerGoroutine = 2_000
		const totalSteps = numGoroutines * stepsPerGoroutine

		counter := NewUpCounter(0)

		results := make([]int64, totalSteps)
		var wg sync.WaitGroup
		var next atomic.Int64

		for i := 0; i != numGoroutines; i++ {

			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j != stepsPerGoroutine; j++ {

					got := counter.Step()
					slot := next.Add(1) - 1
					results[slot] = got
				}
			}()
		}

		wg.Wait()

		require.Equal(t, int64(totalSteps), next.Load())
		require.Equal(t, int64(totalSteps), counter.Load())

		seen := make(map[int64]struct{}, totalSteps)
		for _, v := range results {
			_, exists := seen[v]
			require.False(t, exists, "duplicate Step() return %d", v)
			seen[v] = struct{}{}
		}

		require.Equal(t, totalSteps, len(seen))

		for expect := int64(1); expect <= int64(totalSteps); expect++ {
			_, ok := seen[expect]
			require.True(t, ok, "missing Step() return %d", expect)
		}
	})

	t.Run("concurrent Load() interleaved with Step()", func(t *testing.T) {

		const numSteppers = 8
		const numLoaders = 8
		const stepsPerStepper = 5_000
		const totalSteps = numSteppers * stepsPerStepper

		counter := NewUpCounter(0)

		var steppersWG sync.WaitGroup
		var loadersWG sync.WaitGroup
		var steppersDone atomic.Bool

		for i := 0; i != numSteppers; i++ {

			steppersWG.Add(1)
			go func() {
				defer steppersWG.Done()

				for j := 0; j != stepsPerStepper; j++ {
					_ = counter.Step()
				}
			}()
		}

		for i := 0; i != numLoaders; i++ {

			loadersWG.Add(1)
			go func() {
				defer loadersWG.Done()

				for !steppersDone.Load() {
					_ = counter.Load()
					runtime.Gosched()
				}
			}()
		}

		steppersWG.Wait()
		steppersDone.Store(true)
		loadersWG.Wait()

		require.Equal(t, int64(totalSteps), counter.Load())
	})
}
