package ebigcache_test

import (
	"context"
	"fmt"

	"github.com/ericluj/ebigcache"
)

func Example() {
	cache, _ := ebigcache.New(context.Background(), ebigcache.DefaultConfig())

	cache.Set("my-unique-key", []byte("value"))

	entry, _ := cache.Get("my-unique-key")
	fmt.Println(string(entry))
	// Output: value
}
