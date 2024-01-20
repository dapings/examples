package fast

// FNV-1a 的 Hash 实现，来源 BigCache
// See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function
// See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash

type fnv64a struct{}

// newDefaultHasher returns a new 64-bit FNV-1a Hasher which makes no memory allocations.
// Its Sum64 method will lay the value out in big-endian byte order.
func newDefaultHasher() fnv64a {
	return fnv64a{}
}

const (
	// offset64 FNVa offset basis.
	offset64 = 14695981039346656037
	// prime64 FNVa prime value.
	prime64 = 1099511628211
)

// Sum64 gets the string and returns its uint64 hash value.
func (f fnv64a) Sum64(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash
}
