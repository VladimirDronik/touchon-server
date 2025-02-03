// Простая реализация очереди с приоритетом.

package priority_queue

import "github.com/pkg/errors"

// New создает очередь с приоритетом
// capability - емкость каждого канала с сообщениями
// priorities - кол-во приоритетов (например, 10 - будет 10 приоритетов от 1 до 10)
func New[T any](capability, priorities int) (*PriorityQueue[T], error) {
	switch {
	case capability < 10:
		return nil, errors.Wrap(errors.New("capability < 10"), "priority_queue.New")
	case priorities < 1 || priorities > 50:
		return nil, errors.Wrap(errors.New("priorities < 1 or > 50"), "priority_queue.New")
	}

	o := &PriorityQueue[T]{
		channels: make([]chan T, 0, priorities),
	}

	// Создаем каналы для каждого приоритета
	for i := 0; i < priorities; i++ {
		o.channels = append(o.channels, make(chan T, capability))
	}

	return o, nil
}

type PriorityQueue[T any] struct {
	channels []chan T
}

// Push помещает сообщение в очередь.
func (o *PriorityQueue[T]) Push(v T, priority int) error {
	if priority < 1 || priority > len(o.channels) {
		return errors.Wrap(errors.Errorf("priority < 1 or > %d", len(o.channels)), "PriorityQueue.Push")
	}

	select {
	case o.channels[priority-1] <- v:
		return nil
	default:
		return errors.Wrap(errors.Errorf("канал с сообщениями с приоритетом %d заполнен", priority), "PriorityQueue.Push")
	}
}

// Pop возвращает сообщение из очереди.
// Обходит все каналы по порядку убывания приоритета.
// Если предыдущие очереди пустые и в текущей есть сообщения,
// возвращает сообщение из текущей очереди.
// Выполняется в строгом порядке.
// Есть минус данного алгоритма - если, например, первая
// очередь всегда будет иметь сообщения (будет наполняться
// быстрее, чем опустошаться), то все сообщения из всех
// остальных очередей никогда не обработаются.
// Альтернативный алгоритм может подразумевать
// "разбавление" приоритетных сообщений сообщениями из
// менее приоритетных очередей в некотором соотношении.
func (o *PriorityQueue[T]) Pop() (T, bool) {
	for _, ch := range o.channels {
		select {
		case v, ok := <-ch:
			if ok {
				return v, ok
			}
		default:
		}
	}

	var null T
	return null, false
}
