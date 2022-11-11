package ebigcache

import "time"

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func now() int64 {
	return time.Now().Unix()
}
