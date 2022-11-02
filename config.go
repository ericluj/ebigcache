package ebigcache

type Config struct {
	Shards             int
	Hasher             Hasher
	MaxEntriesInWindow int
	MaxEntrySize       int
}

func DefaultConfig() Config {
	return Config{
		Shards:             1024,
		Hasher:             newDefaultHasher(),
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
	}
}

// shard能保存的entry数量
func (c Config) initialShardSize() int {
	return max(c.MaxEntriesInWindow/c.Shards, minimumEntriesInShard)
}
