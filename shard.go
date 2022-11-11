package ebigcache

import (
	"errors"
	"sync"

	"github.com/ericluj/ebigcache/queue"
)

var (
	ErrEntryNotFound = errors.New("Entry not found")
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
	currentTimestamp := uint64(now())

	s.lock.Lock()

	w := wrapEntry(currentTimestamp, hashedKey, key, entry, &s.entryBuffer)

	index, err := s.entries.Push(w)
	if err != nil {
		s.lock.Unlock()
		return err
	}

	s.hashmap[hashedKey] = uint32(index)
	s.lock.Unlock()
	return nil

}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	s.lock.RLock()

	wrappedEntry, err := s.getWrappedEntry(hashedKey)
	if err != nil {
		s.lock.RUnlock()
		return nil, err
	}

	if entryKey := readKeyFromEntry(wrappedEntry); key != entryKey {
		s.lock.RUnlock()
		return nil, ErrEntryNotFound
	}

	entry := readEntry(wrappedEntry)
	s.lock.RUnlock()

	return entry, nil
}

func (s *cacheShard) getWrappedEntry(hashedKey uint64) ([]byte, error) {
	itemIndex := s.hashmap[hashedKey]

	wrappedEntry, err := s.entries.Get(int(itemIndex))
	if err != nil {
		return nil, err
	}

	return wrappedEntry, err
}
