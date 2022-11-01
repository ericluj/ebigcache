package ebigcache

import "time"

type Clock interface {
	Epoch() int64
}

type systemClock struct {
}

func (c *systemClock) Epoch() int64 {
	return time.Now().Unix()
}
