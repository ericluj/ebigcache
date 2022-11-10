package ebigcache

const (
	minimumEntriesInShard = 10 // Minimum number of entries in single shard
)

type Config struct {
	Hasher              Hasher
	Shards              int
	MaxEntriesInWindow  int // 缓存中最大entry数量
	MaxEntrySizeInBytes int // entry最大字节
}

func DefaultConfig() Config {
	return Config{
		Hasher:              newDefaultHasher(),
		Shards:              1024,
		MaxEntriesInWindow:  1000 * 10 * 60,
		MaxEntrySizeInBytes: 500,
	}
}

// shard中的entry数量
func (c Config) entriesInShard() int {
	return max(c.MaxEntriesInWindow/c.Shards, minimumEntriesInShard)
}
