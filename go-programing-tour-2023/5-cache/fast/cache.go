package fast

type fast struct {
	shards    []*shard
	shardMask uint64
	hash      fnv64a
}

func (f *fast) getShard(key string) *shard {
	hashedKey := f.hash.Sum64(key)
	return f.shards[hashedKey&f.shardMask]
}

func (f *fast) Set(key string, value interface{}) {
	f.getShard(key).set(key, value)
}

func (f *fast) Get(key string) interface{} {
	return f.getShard(key).get(key)
}

func (f *fast) Del(key string) {
	f.getShard(key).del(key)
}

func (f *fast) DelOldest() {
	panic("no implements")
}

func (f *fast) Len() int {
	length := 0
	for _, s := range f.shards {
		length += s.len()
	}

	return length
}
