package ebigcache

type Config struct {
	Shards int
	Hasher Hasher
}

func DefaultConfig() Config {
	return Config{
		Shards: 1024,
		Hasher: newDefaultHasher(),
	}
}
