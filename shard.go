package ebigcache

import (
	"sync"

	"github.com/ericluj/ebigcache/queue"
)

type cacheShard struct {
	hashmap     map[uint64]uint32
	lock        sync.RWMutex
	clock       Clock
	entries     *queue.BytesQueue
	entryBuffer []byte
}

func initNewShard(config Config, clock Clock) *cacheShard {
	// 保存entries的总byte长度
	bytesQueueInitialCapacity := config.initialShardSize() * config.MaxEntrySize
	return &cacheShard{
		// TODO:
		hashmap:     make(map[uint64]uint32, config.initialShardSize()),
		entries:     queue.NewBytesQueue(bytesQueueInitialCapacity),
		entryBuffer: make([]byte, config.MaxEntrySize+headersSizeInBytes),
		clock:       clock,
	}
}

func (s *cacheShard) set(key string, hashedKey uint64, entry []byte) error {
	currentTimestamp := uint64(s.clock.Epoch())

	s.lock.Lock()

	// 当map没有取到值的时候，返回的是类型的零值
	if previousIndex := s.hashmap[hashedKey]; previousIndex != 0 {
		if previousEntry, err := s.entries.Get(int(previousIndex)); err == nil {

		}
	}

	return nil
}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	return []byte("value"), nil
}
