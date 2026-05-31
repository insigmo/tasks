package concurrent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestPatternsFixed(t *testing.T) {
	task2Fixed()
}

// =============================== Task ==============================
// 1. Необходимо провести ревью и рассказать за сколько отработает код
// 2. Необходимо распараллелить, чтобы еще если вдруг упадет, то все горутины закрыть
type User struct {
	Name string
}

func fetch(_ context.Context, user User) (string, error) {
	time.Sleep(time.Millisecond * 10)

	return user.Name, nil
}

func process(ctx context.Context, users []User) (map[string]int64, error) {
	names := make(map[string]int64, 0)

	for _, u := range users {
		name, err := fetch(ctx, u)
		if err != nil {
		}

		names[name] = names[name] + 1
	}

	return names, nil
}

func task2() {
	names := []User{
		{"Ann"},
		{"Bob"},
		{"Cindy"},
		{"Bob"},
	}

	ctx := context.Background()

	start := time.Now()
	res, err := process(ctx, names)
	if err != nil {
		fmt.Println("an error occured:", err.Error())
	}
	fmt.Println("time:", time.Since(start))
	fmt.Println(res)
}

// ============================ Resolution ===========================
// Если посмотреть код, то видим, что в мапе 4 элемента, обработка(функция fetch) на каждый элемент идет примерно 10мс.
// Получается что код будет выполняться примерно 40мс. Нам нужно будет распараллелить, да еще так, чтобы если вдруг упадет
// 1 горутина, другие тоже завершались бы. Для этих целей идеально подходит errgroup. Это либа которая совмещает работу
// с ошибками и WaitGroup. Без errgroup пришлось бы везде раскидывать контексты, код удлинился бы в 2-3 раза.

// Когда допустим получится ошибка в какой-то горутине, то мы функция будет продолжать работать
// еще около 10мс. 10мс это просто имитация полезной нагрузки, оно может быть больше, отправлять запросы и тд
// По факту это особо не мешает, но тем, не менее, если на интервью захочется выделиться, лучше тоже завершать
// функцию досрочно при ошибке. Для этого придется задействовать контекст.
func fetch2(ctx context.Context, user User) (string, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		time.Sleep(time.Millisecond * 10)
	}()

	select {
	case <-done:
		return user.Name, nil
	case <-ctx.Done():
		return "", errors.New("context canceled")
	}
}

func process2(ctx context.Context, users []User) (map[string]int64, error) {
	names := make(map[string]int64)
	// создаем мьютекс, чтобы защитить словарь от состояния гонки
	mw := sync.Mutex{}

	egroup, ectx := errgroup.WithContext(ctx)

	// Мы не знаем сколько юзеров придет, если там много будет, то память переполнится и программа упадет
	// Поэтому мы защищаемся лимитом, чтобы одновременно работало 100 горутин
	egroup.SetLimit(100)

	for _, u := range users {
		// Функция Go запускает в горутине функцию, внутри своей WaitGroup добавляет 1 и также по завершению отнимает
		egroup.Go(func() error {
			name, err := fetch2(ectx, u)
			if err != nil {
				return err
			}

			mw.Lock()
			defer mw.Unlock()
			names[name] = names[name] + 1

			return nil
		})
	}

	// Обязательно ждем пока все горутины не завершатся
	if err := egroup.Wait(); err != nil {
		return nil, err
	}
	return names, nil
}

func task2Fixed() {
	names := []User{
		{"Ann"},
		{"Bob"},
		{"Cindy"},
		{"Bob"},
	}

	ctx := context.Background()

	start := time.Now()
	res, err := process2(ctx, names)
	if err != nil {
		fmt.Println("an error occurred:", err.Error())
	}
	fmt.Println("time:", time.Since(start))
	fmt.Println(res)
}
