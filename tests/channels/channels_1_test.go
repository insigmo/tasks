package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestGoroutine1(t *testing.T) {
	task1Fixed()
}

func worker1() <-chan int {
	ch := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		close(ch)
	}()

	return ch
}

func task1() {
	start := time.Now()
	_, _ = worker1(), worker1()

	fmt.Println(time.Since(start))
}

func task1Fixed() {
	start := time.Now()
	_, _ = <-worker1(), <-worker1() // добавляем чтение из канала, если не добавить, горутины не успеют завершится

	fmt.Println(time.Since(start)) // так как мы последовательно читаем из каналов, приходится ждеть 2 секунды
}
