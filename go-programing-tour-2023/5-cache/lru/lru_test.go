package lru

import (
	"testing"

	"github.com/matryer/is"
)

func TestSet(t *testing.T) {
	th := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	th.Equal(v, 1)

	cache.Del("k1")
	th.Equal(0, cache.Len())

	// cache.Set("k2", time.Now().Local())
}

func TestOnEvicted(t *testing.T) {
	th := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, val interface{}) {
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k2", "k3"}

	th.Equal(expected, keys)
	th.Equal(2, cache.Len())
}
