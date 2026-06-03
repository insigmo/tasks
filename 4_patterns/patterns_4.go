package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// =============================== Task ==============================
func mainTask() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}

	process(urls)
}

// реализовать параллельные запросы по адресам из списка
// подсчитать количество для каждого StatusCode ответа
// предусмотреть возможность отмены запроса по таймауту
//func process(urls []string) {
//}

// ============================ Resolution ===========================
// В первую очередь нужно будет спросить нужно ли с лимитером делать или нет
// Если нужно, уточни с лимитером чего? RPS, горутин или коннектов. Если посмотреть файлик patterns_3.go,
// то можно увидеть решение каждого из них и с ними реализовать.
// Например, сказали с лимитом коннектов.

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}

	fmt.Println(process(urls))
}

// Наиболее популярные ответы это 200, 400, 401, 404, 500, поэтому 5, хотя можно с запасом
const baseStatusCodesAmount = 5

var client http.Client

func process(urls []string) map[int]int {
	// чтобы подсчитать статус коды, можно создать словарик,
	// где ключом будет статус код, а значением - количество
	statusCodes := make(map[int]int, baseStatusCodesAmount)
	wg := sync.WaitGroup{}
	mw := sync.Mutex{}
	wg.Add(len(urls))
	for _, url := range urls {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			mw.Lock()
			defer mw.Unlock()
			statusCodes[resp.StatusCode]++
		}()
	}
	wg.Wait()
	return statusCodes
}
