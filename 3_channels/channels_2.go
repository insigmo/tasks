//go:build ignore

package main

import (
	"fmt"
)

// =============================== Task ==============================
// Необходимо рассказать что получится при выполнении кода и как пофиксить
func mainTask() {
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

// ============================ Resolution ===========================
func main() {
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
