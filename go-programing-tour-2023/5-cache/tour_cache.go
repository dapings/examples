package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter    Getter
}

func NewTourCache(cache Cache, getter Getter) *TourCache {
	return &TourCache{mainCache: newSafeCache(cache), getter: getter}
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}

func (t *TourCache) Set(key string, val interface{}) {
	if val == nil {
		return
	}

	t.mainCache.set(key, val)
}

func (t *TourCache) Get(Key string) interface{} {
	val := t.mainCache.get(Key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(Key)
		if val == nil {
			return nil
		}

		t.mainCache.set(Key, val)
		return val
	}

	return nil
}
