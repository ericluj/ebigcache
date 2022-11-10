package ebigcache

import (
	"sync"

	"github.com/ericluj/ebigcache/queue"
)

type cacheShard struct {
	lock        sync.RWMutex
	hashmap     map[uint64]uint32
	entryBuffer []byte
	entries     *queue.BytesQueue
}

func newShard(c Config) *cacheShard {
	return &cacheShard{
		hashmap:     make(map[uint64]uint32),
		entryBuffer: make([]byte, headersSizeInBytes+c.MaxEntrySizeInBytes),
		entries:     queue.NewBytesQueue(c.entriesInShard()),
	}
}

func (s *cacheShard) set(key string, hashedKey uint64, entry []byte) error {

}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {

}
