package lru

import (
	"container/list"

	cache "github.com/dapings/examples/go-programing-tour-2023/5-cache"
)

// lru 是一个LRU cache，不是并发安全的。
type lru struct {
	// 缓存最大的容量，单位字节；
	// group cache 使用的是最大存放entry个数
	maxBytes int

	// 当一个entry从缓存中移除时，调用此回调函数，默认nil
	// group cache 中的key是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key不算
	usedBytes int

	ll    *list.List
	cache map[string]*list.Element
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (l *lru) Set(key string, value interface{}) {
	// TODO implement me
	panic("implement me")
}

func (l *lru) Get(Key string) interface{} {
	// TODO implement me
	panic("implement me")
}

func (l *lru) Del(key string) {
	// TODO implement me
	panic("implement me")
}

func (l *lru) DelOldest() {
	// TODO implement me
	panic("implement me")
}

func (l *lru) Len() int {
	// TODO implement me
	panic("implement me")
}
