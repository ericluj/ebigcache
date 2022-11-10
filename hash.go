package ebigcache

type Hasher interface {
	Sum64(string) uint64
}

func newDefaultHasher() Hasher {
	return fnv64a{}
}
