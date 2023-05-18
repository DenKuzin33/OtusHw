package hw04lrucache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront("10", 10) // [10]
		l.PushBack("20", 20)  // [10, 20]
		l.PushBack("30", 30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(Key(strconv.Itoa(v)), v)
			} else {
				l.PushBack(Key(strconv.Itoa(v)), v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("front and back equality", func(t *testing.T) {
		l := NewList()
		l.PushBack("1", 1)
		require.NotNil(t, l.Front())

		l.MoveToFront(l.Back())
		require.Equal(t, l.Front(), l.Back())

		l.Remove(l.Front())
		require.Empty(t, l.Len())
		require.Nil(t, l.Back())

		l.PushFront("2", 2)
		l.PushFront("3", 3)
		require.Equal(t, l.Back().Value, 2)

		l.Remove(l.Front().Next)
		require.Equal(t, l.Back(), l.Front())
	})
}
