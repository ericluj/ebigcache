package queue

type BytesQueue struct {
	array    []byte
	capacity int
	head     int
	tail     int
	count    int
}

func NewBytesQueue(capacity int) *BytesQueue {
	return &BytesQueue{
		array:    make([]byte, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		count:    0,
	}
}
