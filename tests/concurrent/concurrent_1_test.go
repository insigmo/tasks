package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrent1(t *testing.T) {
	const maxCount = 100
	var wg sync.WaitGroup

	wg.Add(maxCount)
	for i := 0; i < maxCount; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}

func task1() {
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
