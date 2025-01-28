package entries

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKey_RawKeyCompare(t *testing.T) {
	key1 := NewKey([]byte("key1"), 1)
	key2 := NewKey([]byte("key2"), 2)

	if key1.RawKeyCompare(key2) >= 0 {
		t.Errorf("key1 should be less than key2")
	}

	if key2.RawKeyCompare(key1) <= 0 {
		t.Errorf("key2 should be greater than key1")
	}
}

func TestKey_KeySize(t *testing.T) {
	key := NewKey([]byte("key1"), 1)
	assert.Equal(t, 12, key.EncodedKeySizeInBytes())
}

func TestKey_CompareKeysWithDecendingTimestamp(t *testing.T) {
	key := NewKey([]byte("consensus"), 10)
	assert.Equal(t, -1, key.CompareKeysWithDescendingTimestamp(NewKey([]byte("distributed"), 10)))
	assert.Equal(t, -1, key.CompareKeysWithDescendingTimestamp(NewKey([]byte("consensus"), 5)))
	assert.Equal(t, 0, key.CompareKeysWithDescendingTimestamp(NewKey([]byte("consensus"), 10)))
	assert.Equal(t, 1, key.CompareKeysWithDescendingTimestamp(NewKey([]byte("consensus"), 15)))
	assert.Equal(t, 1, key.CompareKeysWithDescendingTimestamp(NewKey([]byte("accurate"), 10)))
}

func TestEncodedBytes(t *testing.T) {
	key := NewKey([]byte("key1"), 1)
	encodedBytes := key.EncodedBytes()
	assert.Equal(t, 12, len(encodedBytes))
	key2 := DecodeFrom(encodedBytes)
	assert.Equal(t, 0, key.CompareKeysWithDescendingTimestamp(key2))
}
