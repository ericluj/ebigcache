package queue

import "encoding/binary"

type BytesQueue struct {
	array        []byte
	capacity     int
	head         int
	tail         int
	count        int
	headerBuffer []byte
}

func NewBytesQueue(capacity int) *BytesQueue {
	return &BytesQueue{
		array:        make([]byte, capacity),
		capacity:     capacity,
		head:         0,
		tail:         0,
		count:        0,
		headerBuffer: make([]byte, binary.MaxVarintLen32),
	}
}

// 获取到的是entry首部的尾部的字节长度
func getNeededSize(length int) int {
	// header指的是varint编码后数字的字节长度
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
	// varint + data字节长度
	return length + header
}

func (q *BytesQueue) Push(data []byte) (int, error) {
	neededSize := getNeededSize(len(data))

	// push会移动tail指针，所以先保存起来
	index := q.tail

	q.push(data, neededSize)

	return index, nil
}

func (q *BytesQueue) push(data []byte, len int) {
	// headerEntrySize (varint字节长度)
	headerEntrySize := binary.PutUvarint(q.headerBuffer, uint64(len))
	q.copy(q.headerBuffer, headerEntrySize)

	// len - headerEntrySize (data的字节长度)
	q.copy(data, len-headerEntrySize)

	q.count++
}

// 将data追加到q.array
func (q *BytesQueue) copy(data []byte, len int) {
	q.tail += copy(q.array[q.tail:], data[:len])
}

func (q *BytesQueue) Get(index int) ([]byte, error) {
	data, _, err := q.peek(index)
	return data, err
}

// uvarint + data
func (q *BytesQueue) peek(index int) ([]byte, int, error) {
	// blockSize (varint + data字节长度)
	// n (varint字节长度)
	blockSize, n := binary.Uvarint(q.array[index:])

	// index + n (data的首部指针位置)
	// index + blockSize (data的尾部指针位置)
	return q.array[index+n : index+int(blockSize)], int(blockSize), nil
}
