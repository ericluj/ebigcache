package ebigcache

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ericluj/ebigcache/queue"
)

var (
	// ErrEntryNotFound is an error type struct which is returned when entry was not found for provided key
	ErrEntryNotFound = errors.New("Entry not found")
)

type cacheShard struct {
	hashmap     map[uint64]uint32
	lock        sync.RWMutex
	clock       Clock
	entries     *queue.BytesQueue
	entryBuffer []byte
	stats       Stats
	isVerbose   bool
	logger      Logger
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

func (s *cacheShard) removeOldestEntry(reason RemoveReason) error {
	return nil
}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	s.lock.RLock()
	wrapedEntry, err := s.getWrappedEntry(hashedKey)
	if err != nil {
		s.lock.RUnlock()
		return nil, err
	}
	if entryKey := readKeyFromEntry(wrapedEntry); key != entryKey {
		s.lock.RUnlock()
		s.collision()
		if s.isVerbose {
			s.logger.Printf("Collision detected. Both %q and %q have the same hash %x", key, entryKey, hashedKey)
		}
		return nil, ErrEntryNotFound
	}
	entry := readEntry(wrapedEntry)
	s.lock.RUnlock()
	s.hit(hashedKey)

	return entry, nil
}

func (s *cacheShard) getWrappedEntry(hashedKey uint64) ([]byte, error) {
	itemIndex := s.hashmap[hashedKey]

	if itemIndex == 0 {
		s.miss()
		return nil, ErrEntryNotFound
	}

	wrappedEntry, err := s.entries.Get(int(itemIndex))
	if err != nil {
		s.miss()
		return nil, err
	}

	return wrappedEntry, err
}

func (s *cacheShard) miss() {
	atomic.AddInt64(&s.stats.Misses, 1)
}

func (s *cacheShard) collision() {
	atomic.AddInt64(&s.stats.Collisions, 1)
}

func (s *cacheShard) hit(key uint64) {
	atomic.AddInt64(&s.stats.Hits, 1)
	if s.statsEnabled {
		s.lock.Lock()
		s.hashmapStats[key]++
		s.lock.Unlock()
	}
}
