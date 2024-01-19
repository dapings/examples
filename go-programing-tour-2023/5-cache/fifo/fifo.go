package fifo

import (
	"container/list"

	cache "github.com/dapings/examples/go-programing-tour-2023/5-cache"
)

// fifo 是一个 FIFO cache。它不是并发安全的。
type fifo struct {
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

type entry struct {
	key string
	val interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.val)
}

// New 创建一个新的 Cache，如果 maxBytes 是 0，表示没有容量限制
func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// Set 往 Cache 尾部增加一个元素（如果已经存在，则放入尾部，并更新值）
func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.val) + cache.CalcLen(value)
		en.val = value

		return
	}

	en := &entry{
		key: key,
		val: value,
	}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

// Get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (f *fifo) Get(Key string) interface{} {
	if e, ok := f.cache[Key]; ok {
		// f.ll.MoveToBack(e)
		return e.Value.(*entry).val
	}

	return nil
}

// Del 从 cache 中删除 key 对应的元素
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

// DelOldest 从 cache 中删除最旧的记录
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

// Len 返回当前 cache 中的记录数
func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.val)
	}
}
