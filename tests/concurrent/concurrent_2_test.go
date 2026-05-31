package concurrent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrent2WithMutex(t *testing.T) {
	const maxCount = 100
	var wg sync.WaitGroup
	wg.Add(maxCount)
	lock := sync.Mutex{}

	counter := 0
	for i := 0; i < maxCount; i++ {
		go func() {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			counter++

		}()
	}
	wg.Wait()

	t.Run("TestConcurrent2WithMutex", func(t *testing.T) {
		if counter != maxCount {
			t.Fail()
		}
	})

}
func TestConcurrent2WithAtomic(t *testing.T) {
	const maxCount = 100
	var wg sync.WaitGroup
	wg.Add(maxCount)

	var counter int64 = 0
	for i := 0; i < maxCount; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	t.Run("TestConcurrent2WithAtomic", func(t *testing.T) {
		if counter != maxCount {
			t.Fail()
		}
	})
}
func task2() {
	// Как исправить?
	counter := 0
	for i := 0; i < 100; i++ {
		go func() {
			counter++
		}()
	}

	fmt.Println(counter)
}
