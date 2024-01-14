package cache

import (
	"log"
	"sync"
	"testing"

	"github.com/dapings/examples/go-programing-tour-2023/5-cache/lru"
	"github.com/matryer/is"
)

func TestTourCacheGet(t *testing.T) {
	data := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}

	getter := GetFunc(func(key string) interface{} {
		log.Println("from data find key", key)

		if val, ok := data[key]; ok {
			return val
		}

		return nil
	})

	tourCache := NewTourCache(lru.New(0, nil), getter)

	th := is.New(t)

	var wg sync.WaitGroup
	for k, v := range data {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()

			th.Equal(tourCache.Get(k), v)

			th.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	th.Equal(tourCache.Get("unknown"), nil)
	th.Equal(tourCache.Get("unknown"), nil)

	th.Equal(tourCache.Stat().NGet, 10)
	th.Equal(tourCache.Stat().NHit, 4)
}
