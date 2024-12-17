package entries

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const TimestampSize = int(unsafe.Sizeof(uint64(0)))

type Key struct {
	Key       []byte
	timestamp uint64
}

var EmptyKey = Key{Key: []byte{}}

func DecodeFrom(buffer []byte) Key {
	if len(buffer) < TimestampSize {
		panic("buffer is too small")
	}

	length := len(buffer)
	key := buffer[:length-TimestampSize]
	timestamp := binary.LittleEndian.Uint64(buffer[length-TimestampSize:])
	return Key{
		Key:       key,
		timestamp: timestamp,
	}
}

func NewKey(key []byte, timestamp uint64) Key {
	return Key{Key: key, timestamp: timestamp}
}

func (k Key) IsLessThanOrEqualTo(key Key) bool {
	cmp := bytes.Compare(k.Key, key.Key)
	if cmp == 0 {
		return k.timestamp <= key.timestamp
	}
	return cmp < 0
}

func (k Key) EqualTo(key Key) bool {
	return bytes.Compare(k.Key, key.Key) == 0 && k.timestamp == key.timestamp
}

func (k Key) Compare(key Key) int {
	return bytes.Compare(k.Key, key.Key)
}

func (k Key) String() string {
	return string(k.Key)
}

func (k Key) CompareKeysWithDecendingTimestamp(otherKey Key) int {
	comparison := bytes.Compare(k.Key, otherKey.Key)
	if comparison != 0 {
		return comparison
	}
	if k.timestamp < otherKey.timestamp {
		return 1
	}
	if k.timestamp > otherKey.timestamp {
		return -1
	}

	return 0
}

func NewStringKey(key string) Key {
	return Key{Key: []byte(key)}
}

func CompareKeys(key1 Key, key2 Key) int {
	return key1.CompareKeysWithDecendingTimestamp(key2)
}

func (k Key) isRawKeyEqual(key Key) bool {
	return bytes.Equal(k.Key, key.Key)
}

func (k Key) RawKeySize() int {
	return len(k.Key)
}

func (k Key) RawKeyCompare(key Key) int {
	return bytes.Compare(k.Key, key.Key)
}

func (k Key) RawKey() []byte {
	return k.Key
}

func (k Key) Timestamp() uint64 {
	return k.timestamp
}

func (k Key) EncodedKeySizeInBytes() int {
	if k.RawKeySize() == 0 {
		return 0
	}
	return len(k.RawKey()) + TimestampSize
}

func (k Key) EncodedBytes() []byte {
	if k.RawKeySize() == 0 {
		return nil
	}
	buffer := make([]byte, k.EncodedKeySizeInBytes())
	copy(buffer, k.RawKey())
	binary.LittleEndian.PutUint64(buffer[k.RawKeySize():], k.Timestamp())
	return buffer
}
