package concurrent

import (
	"fmt"
	"testing"
)

func TestGoroutine3(t *testing.T) {
	task3Fixed1()
	//task3Fixed2()
}

func spawnMessages(n int) chan string {
	ch := make(chan string, 1)

	for i := 0; i < n; i++ {
		ch <- fmt.Sprintf("msg %d", i+1)
	}
	return ch
}

func task3() {
	n := 10

	for msg := range spawnMessages(n) {
		fmt.Println("received:", msg)
	}
}

// Первый способ(рекомендуется). Так как мы читаем из канала в цикле, то канал который возвращается
// должен содержать данные, но так получается что канал у нас хоть и буферизированный,
// размер буфера всего лишь 1.
// На интервью нужно будет сказать что после того как в канал поместится "msg 0",
// второй раз записать не получится, потому что с канала еще никто не читает.
// Мы же еще не вернули канал ведь.
// Поэтому нужно будет в отдельной горутине запустить цикл, который будет снабжать канал данными,
// после чего закрыть обязательно канал в этой же горутине по завершению.
func spawnMessagesFixed1(n int) chan string {
	ch := make(chan string, 1)

	go func() { // была добавлена отдельная горутина
		defer close(ch) // канал закрыть внутри горутины(обязательно внутри горутины)
		for i := 0; i < n; i++ {
			ch <- fmt.Sprintf("msg %d", i+1)
		}
	}()
	return ch
}

func task3Fixed1() {
	n := 10

	for msg := range spawnMessagesFixed1(n) {
		fmt.Println("received:", msg)
	}
}

// Второй способ(менее рекомендуемый). Так как мы читаем из канала в цикле, то канал который
// возвращается должен содержать данные, но так получается что канал у нас хоть и буферизированный,
// размер буфера всего лишь 1.
// Можно поднять размер буфера на n и потом закрыть канал.
// Вообще эту задачу можно решить и не буферизированным каналом,
// просто вынести запись в отдельную горутину, но в данной задаче нужно решить с буфером
func spawnMessagesFixed2(n int) chan string {
	ch := make(chan string, n) // был изменен размер на n
	defer close(ch)
	for i := 0; i < n; i++ {
		ch <- fmt.Sprintf("msg %d", i+1)
	}
	return ch
}

func task3Fixed2() {
	n := 10

	for msg := range spawnMessagesFixed1(n) {
		fmt.Println("received:", msg)
	}
}
