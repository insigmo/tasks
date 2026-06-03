//go:build ignore

package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	//task3FixedWithConnectLimiter()
	//task3FixedWithGoroutineLimiter()
	task3FixedWithRpsLimiter()
}

// =============================== Task ==============================
// Задача заключается в том, чтобы отправить все запросы, но сделать это через рейт-лимитер, то есть ограничить количество
// одновременных запросов. Это делается чаще всего чтобы не заддосить сервер, куда отправляется запрос.

type Request struct {
	Payload string
}

type Client interface {
	SendRequest(ctx context.Context, request Request) error
	WithLimiter(ctx context.Context, requests []Request)
}

type client struct {
}

func (c client) SendRequest(ctx context.Context, request Request) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("sending request", request.Payload)
	return nil
}

func (c client) WithLimiter(ctx context.Context, requests []Request) {
}

func mainTask() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.WithLimiter(ctx, requests)
}

// ============================ Resolution ===========================
// Перед тем как приступить к решению задач, нужно будет спросить, что именно нужно будет ограничить.
// Бывает 3 способа ограничения:
// 1. Ограничение коннектов
// 2. Ограничение горутин
// 3. Ограничение запросов в секунду(rps).

// Поэтому тут будет 3 решения.
// 1. Ограничение коннектов, допустим 10 одновременных коннектов

const maxConnect = 10

func (c client) WithConnectLimiter(ctx context.Context, ch <-chan Request) {
	wg := sync.WaitGroup{}
	wg.Add(maxConnect)

	for range maxConnect {
		go func() {
			defer wg.Done()
			for req := range ch {
				if err := c.SendRequest(ctx, req); err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
	}
	wg.Wait()
}

// Лучше всего читать данные, в данном случае реквести через канал.
// Поэтому и создаем генератор, который нам вернет канал реквестов.
func generate(reqs []Request) <-chan Request {
	out := make(chan Request)

	go func() {
		defer close(out)

		for _, req := range reqs {
			out <- req
		}
	}()

	return out
}

func task3FixedWithConnectLimiter() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.WithConnectLimiter(ctx, generate(requests))
}

// 2. Ограничение горутин, допустим 100 горутин будут работать одновременно.
// Легче всего опять это сделать через каналы.
const maxGoroutines = 100

func (c client) task3FixedWithGoroutineLimiter(ctx context.Context, requests []Request) {
	threads := make(chan struct{}, maxGoroutines)
	go func() {
		for range maxGoroutines {
			threads <- struct{}{}
		}
	}()

	for _, req := range requests {
		go func() {
			go func() {
				<-threads
			}()
			if err := c.SendRequest(ctx, req); err != nil {
				fmt.Println(err)
				return
			}
			threads <- struct{}{}

		}()
	}
	for range maxGoroutines {
		<-threads
	}
}

func task3FixedWithGoroutineLimiter() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.task3FixedWithGoroutineLimiter(ctx, requests)
}

// 3. Ограничение запросов в секунду(rps), допустим 100 rps.
// Легче всего опять это сделать через тикер
const rps = 100

func (c client) task3FixedWithRpsLimiter(ctx context.Context, requests []Request) {
	ticker := time.NewTicker(time.Second / time.Duration(rps))

	for _, req := range requests {
		<-ticker.C
		go func() {
			if err := c.SendRequest(ctx, req); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}

func task3FixedWithRpsLimiter() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.task3FixedWithRpsLimiter(ctx, requests)
}
