package parallel

import (
	"errors"
	"sync"
	"time"
)

var ErrTimeout = errors.New("timeout exceeded")

type Task func()

// Do запускает задачи параллельно в несколько потоков.
// Ждет либо завершения всех задач, либо наступления таймаута.
func Do(threads int, tasks []Task, timeoutForAllTasks time.Duration) error {
	tasksList := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksList <- task
	}
	close(tasksList)

	wg := &sync.WaitGroup{}
	wg.Add(threads)

	errs := make(chan error, 1)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()

			timer := time.NewTimer(timeoutForAllTasks)
			defer timer.Stop()

			for task := range tasksList {
				select {
				case <-timer.C:
					// Один из потоков сможет записать ошибку таймаута
					select {
					case errs <- ErrTimeout:
					default:
					}

					return
				case <-DoAsync(task):
				}
			}
		}()
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return <-errs
	}

	return nil
}

// DoAsync запускает задачу и возвращает канал,
// который будет закрыт по завершении задачи.
func DoAsync(task Task) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		task()
		close(done)
	}()
	return done
}

// DoWithTimeout запускает задачу и ждет либо завершения
// задачи, либо наступления таймаута.
func DoWithTimeout(task Task, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	done := make(chan struct{})

	go func() {
		task()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timer.C:
		return ErrTimeout
	}
}
