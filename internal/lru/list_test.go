package lru

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("nil does not break move and remove methods", func(t *testing.T) {
		l := NewList()

		assert.NotPanics(t, func() {
			l.MoveToFront(nil)
		})

		assert.NotPanics(t, func() {
			l.Remove(nil)
		})
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront("3") // 3
		l.PushFront("2") // 2,3
		l.PushFront("1") // 1,2,3

		require.Equal(t, 3, l.Len())

		require.Equal(t, "1", l.Front().Value)
		require.Equal(t, "2", l.Front().Next.Value)
		require.Equal(t, "3", l.Front().Next.Next.Value)

		// drop middle
		l.Remove(l.Front().Next)
		require.Equal(t, 2, l.Len())
		require.Equal(t, "1", l.Front().Value)
		require.Equal(t, "3", l.Back().Value)

		//move front
		l.PushBack("4")
		require.Equal(t, "4", l.Back().Value)
		l.MoveToFront(l.Back())
		require.Equal(t, "4", l.Front().Value)
	})
}
