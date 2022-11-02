package queue

type BytesQueue struct {
	capacity int
	array    []byte
}

func NewBytesQueue(capacity int) *BytesQueue {
	return &BytesQueue{
		capacity: capacity,
		array:    make([]byte, capacity),
	}
}

func (q *BytesQueue) Get(index int) ([]byte, error) {
	return nil, nil
}
