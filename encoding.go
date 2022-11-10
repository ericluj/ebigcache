package ebigcache

import "encoding/binary"

const (
	timestampSizeInBytes = 8
	hashSizeInBytes      = 8
	keySizeInBytes       = 2
	headersSizeInBytes   = timestampSizeInBytes + hashSizeInBytes + keySizeInBytes
)

// timestamp(8) + hash(8) + keyLen(2) + key + entry
func wrapEntry(timestamp uint64, hash uint64, key string, entry []byte, buffer *[]byte) []byte {
	keyLength := len(key)
	blobLength := headersSizeInBytes + keyLength + len(entry)

	if blobLength > len(*buffer) {
		*buffer = make([]byte, blobLength)
	}
	blob := *buffer

	binary.LittleEndian.PutUint64(blob, timestamp)
	binary.LittleEndian.PutUint64(blob[timestampSizeInBytes:], hash)
	binary.LittleEndian.PutUint16(blob[timestampSizeInBytes+hashSizeInBytes:], uint16(keyLength))
	copy(blob[headersSizeInBytes:], key)
	copy(blob[headersSizeInBytes+keyLength:], entry)

	return blob[:blobLength]
}

// 从wrappedEntry中读取时间戳
func readTimestampFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

// 从wrappedEntry中读取hashedKey
func readHashFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data[timestampSizeInBytes:])
}

// 从wrappedEntry中读取key
func readKeyFromEntry(data []byte) string {
	keyLength := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, keyLength)
	copy(dst, data[headersSizeInBytes:headersSizeInBytes+keyLength])

	return bytesToString(dst)
}

// 从wrappedEntry中读取entry
func readEntry(data []byte) []byte {
	keyLength := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, len(data)-int(headersSizeInBytes+keyLength))
	copy(dst, data[headersSizeInBytes+keyLength:])

	return dst
}
