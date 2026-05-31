package concurrent

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrent2WithMutex(t *testing.T) {
	task2FixedWithMutex()
	//task2FixedWithAtomic()
}

// =============================== Task ==============================
// Что произойдет если запустить код и что выведется?
func task2() {
	counter := 0
	for i := 0; i < 100; i++ {
		go func() {
			counter++
		}()
	}

	fmt.Println(counter)
}

// ============================ Resolution ===========================
// При запуске функции произойдет состояние гонки(race condition)
// В этом состоянии, несколько горутин могут взять counter равным условно 0, каждый добавит 1 и потом вернет.
// Получается, что 2 или более горутины сделали одно и тоже несколько раз.
// Надо защитить код от состояния гонки, для этого тут можно воспользоваться примитивами синхронизации
// 1 способ это можно воспользоваться мьютексом
func task2FixedWithMutex() {
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
	fmt.Println(counter)

}

// 2 способ. Воспользоваться атомиком. Атомик это либа которая выполняет некоторые операции атомарными,
// то есть, выполняет какое-то действие за 1 процессорный шаг. Например, если мы сделаем 1+1, с точки зрения процессора
// нам нужно создать две единицы, затем их сложить, затем сохранить результат и вернуть его, а атомик это сделает за 1 шаг.
// Поэтому если даже параллельно будет выполняться, то другая горутина не может вклиниться и что-то испортить.
func task2FixedWithAtomic() {
	const maxCount = 100
	var wg sync.WaitGroup
	wg.Add(maxCount)

	var counter int64 = 0
	for i := 0; i < maxCount; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // тут вместо counter++ мы делаем добавление другого числа
		}()
	}
	fmt.Println(counter)
}
