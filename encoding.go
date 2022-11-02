package ebigcache

import "encoding/binary"

const (
	timestampSizeInBytes = 8
	hashSizeInBytes      = 8
	keySizeInBytes       = 2
	headersSizeInBytes   = timestampSizeInBytes + hashSizeInBytes + keySizeInBytes
)

// timestamp(8b) + hash(8b) + keyLen(2b) + key + entry
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

func readKeyFromEntry(data []byte) string {
	keyLength := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, keyLength)
	copy(dst, data[headersSizeInBytes:headersSizeInBytes+keyLength])

	return bytesToString(dst)
}
