package queue

import "encoding/binary"

type BytesQueue struct {
	capacity     int
	array        []byte
	headerBuffer []byte
}

func NewBytesQueue(capacity int) *BytesQueue {
	return &BytesQueue{
		capacity:     capacity,
		array:        make([]byte, capacity),
		headerBuffer: make([]byte, binary.MaxVarintLen32),
	}
}

func (q *BytesQueue) Get(index int) ([]byte, error) {
	return nil, nil
}

// TODO:
func getNeededSize(length int) int {
	var header int
	switch {
	case length < 127: // 1<<7-1
		header = 1
	case length < 16382: // 1<<14-2
		header = 2
	case length < 2097149: // 1<<21 -3
		header = 3
	case length < 268435452: // 1<<28 -4
		header = 4
	default:
		header = 5
	}

	return length + header
}

func (q *BytesQueue) Push(data []byte) (int, error) {
	return 0, nil
}
