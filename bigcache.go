package ebigcache

import "context"

type BigCache struct {
	config    Config
	hash      Hasher
	shards    []*cacheShard
	shardMask uint64
}

func New(ctx context.Context, config Config) (*BigCache, error) {
	return newBigCache(ctx, config)
}

func newBigCache(ctx context.Context, config Config) (*BigCache, error) {
	if config.Hasher == nil {
		config.Hasher = newDefaultHasher()
	}

	cache := &BigCache{
		config:    config,
		hash:      config.Hasher,
		shards:    make([]*cacheShard, config.Shards),
		shardMask: uint64(config.Shards - 1),
	}

	for i := 0; i < config.Shards; i++ {
		cache.shards[i] = newShard(config)
	}

	return cache, nil
}

func (c *BigCache) getShard(hashedKey uint64) (shard *cacheShard) {
	return c.shards[hashedKey&c.shardMask]
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
