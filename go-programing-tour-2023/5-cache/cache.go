package cache

import (
	"sync"
)

// Cache 缓存接口
type Cache interface {
	Set(key string, value interface{})
	Get(Key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

// DefaultMaxBytes 默认允许占用的最大内存
const DefaultMaxBytes = 1 << 29

// safeCache 并发安全缓存
type safeCache struct {
	m     sync.Mutex
	cache Cache

	nget, nhit int
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{cache: cache}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()

	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.Lock()
	defer sc.m.Unlock()

	sc.nget++
	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		sc.nhit++
	}

	return v
}

func (sc *safeCache) stat() *Stat {
	sc.m.Lock()
	defer sc.m.Unlock()

	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}

type Stat struct {
	NHit, NGet int
}
