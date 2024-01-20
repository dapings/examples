package fast

import (
	"container/list"
	"sync"
)

type shard struct {
	locker sync.RWMutex

	// 最大存放 entry 个数
	maxEntries int

	// 当一个entry从缓存中移除时，调用此回调函数，默认nil
	// group cache 中的key是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key string
	val interface{}
}

// 创建一个新的 shard，如果 maxBytes 是 0，表示没有容量限制
func newShard(maxEntries int, onEvicted func(key string, value interface{})) *shard {
	return &shard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// set 往 Cache 尾部增加一个元素（如果已经存在，则放入尾部，并更新值）
func (s *shard) set(key string, val interface{}) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if e, ok := s.cache[key]; ok {
		s.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.val = val
		return
	}

	en := &entry{key: key, val: val}
	e := s.ll.PushBack(en)
	s.cache[key] = e

	if s.maxEntries > 0 && s.ll.Len() > s.maxEntries {
		s.removeElement(s.ll.Front())
	}
}

// get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (s *shard) get(key string) interface{} {
	s.locker.RLocker()
	defer s.locker.RUnlock()

	if e, ok := s.cache[key]; ok {
		s.ll.MoveToBack(e)
		return e.Value.(*entry).val
	}

	return nil
}

// del 从 cache 中删除 key 对应的元素
func (s *shard) del(key string) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if e, ok := s.cache[key]; ok {
		s.removeElement(e)
	}
}

// delOldest 从 cache 中删除最旧的记录
func (s *shard) delOldest() {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.removeElement(s.ll.Front())
}

// len 返回当前 cache 中的记录数
func (s *shard) len() int {
	s.locker.RLocker()
	defer s.locker.Unlock()

	return s.ll.Len()
}

func (s *shard) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	s.ll.Remove(e)
	en := e.Value.(*entry)
	delete(s.cache, en.key)

	if s.onEvicted != nil {
		s.onEvicted(en.key, en.val)
	}
}
