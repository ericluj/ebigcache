package ebigcache

import "sync"

type cacheShard struct {
	hashmap map[uint64]uint32
	lock    sync.RWMutex
	clock   Clock
}

func initNewShard(config Config, clock Clock) *cacheShard {
	return &cacheShard{
		// TODO:
		hashmap: make(map[uint64]uint32, config.Shards),
		clock:   clock,
	}
}

func (s *cacheShard) set(key string, hashedKey uint64, entry []byte) error {
	return nil
}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	return []byte("value"), nil
}
