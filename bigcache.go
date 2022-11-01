package ebigcache

import "context"

type BigCache struct {
	shards []*cacheShard
	hash   Hasher
	config Config
	clock  Clock
}

func New(ctx context.Context, config Config) (*BigCache, error) {
	return newBigCache(ctx, config, &systemClock{})
}

func newBigCache(ctx context.Context, config Config, clock Clock) (*BigCache, error) {
	if config.Hasher == nil {
		config.Hasher = newDefaultHasher()
	}

	cache := &BigCache{
		shards: make([]*cacheShard, config.Shards),
		clock:  clock,
		hash:   config.Hasher,
		config: config,
	}

	for i := 0; i < config.Shards; i++ {
		cache.shards[i] = initNewShard(config, clock)
	}

	return cache, nil
}

func (c *BigCache) getShard(hashedKey uint64) (shard *cacheShard) {
	// TODO:
	return c.shards[hashedKey]
}

func (c *BigCache) Set(key string, entry []byte) error {
	hashedKey := c.hash.Sum64(key)
	shard := c.getShard(hashedKey)
	return shard.set(key, hashedKey, entry)
}

func (c *BigCache) Get(key string) ([]byte, error) {
	hashedKey := c.hash.Sum64(key)
	shard := c.getShard(hashedKey)
	return shard.get(key, hashedKey)
}
