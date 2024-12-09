package entries

import "bytes"

type Key struct {
	Key []byte
}

func NewKey(key []byte) *Key {
	return &Key{Key: key}
}

func (k Key) IsLessThanOrEqual(key Key) bool {
	return bytes.Compare(k.Key, key.Key) <= 0
}

func (k Key) EqualTo(key Key) bool {
	return bytes.Compare(k.Key, key.Key) == 0
}

func (k Key) Compare(key Key) int {
	return bytes.Compare(k.Key, key.Key)
}

func (k Key) String() string {
	return string(k.Key)
}

func NewStringKey(key string) Key {
	return Key{Key: []byte(key)}
}
