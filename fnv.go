package ebigcache

func newDefaultHasher() Hasher {
	return &fnv64a{}
}

type fnv64a struct{}

// TODO:
func (f *fnv64a) Sum64(key string) uint64 {
	return 0
}
