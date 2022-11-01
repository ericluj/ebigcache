package ebigcache

type Hasher interface {
	Sum64(string) uint64
}
