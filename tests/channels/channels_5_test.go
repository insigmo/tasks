package concurrent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGoroutine5(t *testing.T) {
	task5Fixed()
}

// =============================== Task ==============================
// Нужно превратить непредсказуемую функцию в предсказуемую.
// То есть если функция выполняется дольше 5 секунд, сделать так, чтобы он завершался.
// При этом нельзя менять функцию unpredictableFunc

func unpredictableFunc() int {
	n := rand.Intn(40)
	time.Sleep(time.Duration(n) * time.Second)
	return n
}

func predictableFunc() int {
	return 0
}

func task5() {
	_ = predictableFunc()
}

// ============================ Resolution ===========================
// Для того, чтобы добавить предсказуемость в контекст, надо добавить тайм-аут, если он не задан
// Если он задан, то через метод Deadline можно проверить, он возвращает время и ок как булевое значение
// Если ок будет ложным, это значит что тайм-аут не задан.
func predictableFuncFixed(ctx context.Context) (int, error) {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

	}
	// нужно также добавить канал, чтобы подождать результат на тот случай, когда непредсказуемый канал
	// завершится раньше тайм-аута, тогда мы закроем канал и селект уже выберет кейс,
	// где канал done завершился и вернет результат
	done := make(chan struct{})
	var res int

	// мы должны запустить непредсказуемую функцию в горутине,
	// чтобы у нас была возможность запустить отдельно какую то логику
	go func() {
		defer close(done)
		res = unpredictableFunc()
	}()

	// C помощью select уже выбираем либо завершенный канал done, либо завершенный канал контекста
	// канал контекста можно получить через метод Done и он завершается когда время прошло
	select {
	case <-done:
		return res, nil
	case <-ctx.Done():
		return 0, errors.New("timed out")
	}
}

// Чтобы поведение было предсказуемым, нужно это поведение создать
// Для этого в Go добавили context. Благодаря чему можно повесить таймауты и тд
// В данном случае можем передать пустой контекст.
func task5Fixed() {
	res, err := predictableFuncFixed(context.Background())
	fmt.Printf("res: %v, err: %v\n", res, err)
}
