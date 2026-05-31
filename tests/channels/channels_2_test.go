package concurrent

import (
	"fmt"
	"testing"
)

func TestGoroutine2(t *testing.T) {
	task2Fixed()
}

func task2() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		fmt.Println(n)
	}
}

func task2Fixed() {
	ch := make(chan int)

	go func() {
		defer close(ch) // необходимо закрыть канал, чтобы цикл при чтении завершился
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		fmt.Println(n)
	}
}
