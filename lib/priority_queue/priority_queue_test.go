package priority_queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New[int](0, 1)
	require.NotNil(t, err, "capability=0 must return error")

	_, err = New[int](10, 0)
	require.NotNil(t, err, "priorities=0 must return error")

	_, err = New[int](10, 51)
	require.NotNil(t, err, "priorities=51 must return error")

	_, err = New[int](10, 1)
	require.Nil(t, err)

	_, err = New[int](100, 50)
	require.Nil(t, err)
}

func TestPriorityQueue_Push(t *testing.T) {
	q, err := New[int](10, 3)
	require.Nil(t, err)

	// Заполняем все каналы
	for i := 0; i < 10; i++ {
		for j := 0; j < 3; j++ {
			err = q.Push(0, j+1)
			require.Nil(t, err)
		}
	}

	err = q.Push(0, 0)
	require.NotNil(t, err)

	err = q.Push(0, 1)
	require.NotNil(t, err)

	err = q.Push(0, 2)
	require.NotNil(t, err)

	err = q.Push(0, 0)
	require.NotNil(t, err)

	err = q.Push(0, 4)
	require.NotNil(t, err)
}

func TestPriorityQueue_Pop(t *testing.T) {
	t.Run("continuously", func(t *testing.T) {
		q, err := New[int](10, 4)
		require.Nil(t, err)

		// Заполняем все каналы
		for i := 1; i <= 10; i++ {
			err = q.Push(i, 3)
			require.Nil(t, err)
		}
		for i := 11; i <= 20; i++ {
			err = q.Push(i, 1)
			require.Nil(t, err)
		}

		for i := 11; i <= 20; i++ {
			v, ok := q.Pop()
			require.Equal(t, true, ok)
			require.Equal(t, i, v)
		}

		for i := 1; i <= 10; i++ {
			v, ok := q.Pop()
			require.Equal(t, true, ok)
			require.Equal(t, i, v)
		}

		v, ok := q.Pop()
		require.Equal(t, false, ok)
		require.Equal(t, 0, v)
	})

	t.Run("parallel", func(t *testing.T) {
		q, err := New[int](10, 4)
		require.Nil(t, err)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// Заполняем все каналы
		go func() {
			for i := 1; i <= 10; i++ {
				err = q.Push(i, 3)
				require.Nil(t, err)
			}
			wg.Done()
		}()
		go func() {
			for i := 11; i <= 20; i++ {
				err = q.Push(i, 1)
				require.Nil(t, err)
			}
			wg.Done()
		}()
		wg.Wait()

		for i := 11; i <= 20; i++ {
			v, ok := q.Pop()
			require.Equal(t, true, ok)
			require.Equal(t, i, v)
		}

		for i := 1; i <= 10; i++ {
			v, ok := q.Pop()
			require.Equal(t, true, ok)
			require.Equal(t, i, v)
		}

		v, ok := q.Pop()
		require.Equal(t, false, ok)
		require.Equal(t, 0, v)
	})
}
