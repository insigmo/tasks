package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestGoroutine1(t *testing.T) {
	task1Fixed()
}

// =============================== Task ==============================
// Необходимо рассказать что получится при выполнении кода
func worker() <-chan int {
	ch := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		close(ch)
	}()

	return ch
}

func task1() {
	start := time.Now()
	_, _ = worker(), worker()

	fmt.Println(time.Since(start))
}

// ============================ Resolution ===========================
func task1Fixed() {
	start := time.Now()
	_, _ = <-worker(), <-worker() // добавляем чтение из канала, если не добавить, горутины не успеют завершится

	fmt.Println(time.Since(start)) // так как мы последовательно читаем из каналов, приходится ждать 2 секунды
}
