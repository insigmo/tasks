//go:build ignore

package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	close(ch1)

	ch2 := make(chan int, 2)
	ch2 <- 3
	ch2 <- 4
	close(ch2)

	var got []int
	for v := range faninFixed(ch1, ch2) {
		got = append(got, v)
	}

	sort.Ints(got)
	want := []int{1, 2, 3, 4}
	for i, v := range want {
		if got[i] != v {
			fmt.Printf("got %v, want %v", got, want)
		}
	}
}

// =============================== Task ==============================
// Объединить каналы в 1 канал
func fanin(chans ...<-chan int) <-chan int {
	results := make(chan int)

	return results
}

// ============================ Resolution ===========================
// Fan-in это паттерн конкурентного программирования при котором объединяются данные со всех каналов в 1 канал
// Для fan-in чаще всего используют генератор. Как можно заметить в функции ниже,
// создается канал, горутина и возвращается канал. Это и есть паттерн генератор.
// Внутри горутины мы создаем цикл, который генерит горутины также, чтобы при считывании из канала не
// блочилась другая горутина. Но нам также нужно закрыть результирующий канал, ибо на выходе, тот кто будет читать,
// ему нужен сигнал, что данные закончились. Закрытие канала как раз идеально подходит для этого, но тут возникает вопрос.
// Когда нужно закрывать канал? Очевидно что когда все каналы залили свои данные в результирующий канал,
// и только потом закрывать его, поэтому нужно добавить примитив синхронизации. WaitGroup подходит для этого.
func faninFixed(chans ...<-chan int) <-chan int {
	results := make(chan int)
	go func() {
		wg := sync.WaitGroup{}

		for _, ch := range chans {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					results <- v
				}
			}()
		}

		wg.Wait()
		close(results)
	}()

	return results
}
