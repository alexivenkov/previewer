package lru

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("a", "1")
		require.False(t, wasInCache)

		wasInCache = c.Set("b", "2")
		require.False(t, wasInCache)

		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, "1", val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, "2", val)

		wasInCache = c.Set("a", "3")
		require.True(t, wasInCache)

		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, "3", val)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Empty(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(1)

		c.Set("test", "1")
		val, _ := c.Get("test")
		require.Equal(t, val, "1")

		c.Clear()
		_, ok := c.Get("test")
		require.False(t, ok)
	})

}
