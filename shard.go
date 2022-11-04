package ebigcache

import (
	"fmt"
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

	// 如果hash已经有值，删掉（当map没有取到值的时候，返回的是类型的零值）
	if previousIndex := s.hashmap[hashedKey]; previousIndex != 0 {
		if previousEntry, err := s.entries.Get(int(previousIndex)); err == nil {
			resetKeyFromEntry(previousEntry)
			delete(s.hashmap, hashedKey)
		}
	}

	w := wrapEntry(currentTimestamp, hashedKey, key, entry, &s.entryBuffer)

	for {
		if index, err := s.entries.Push(w); err == nil {
			s.hashmap[hashedKey] = uint32(index)
			s.lock.Unlock()
			return nil
		}
		if s.removeOldestEntry(NoSpace) != nil {
			s.lock.Unlock()
			return fmt.Errorf("entry is bigger than max shard size")
		}
	}
}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	return []byte("value"), nil
}

func (s *cacheShard) removeOldestEntry(reason RemoveReason) error {
	return nil
}
