package block

import (
	"github.com/stretchr/testify/assert"
	"lsm/entries"
	"testing"
)

func TestBuilder_Add_Success(t *testing.T) {
	builder := NewBuilder(DefaultBlockSize)
	key := entries.NewKey([]byte("key1"), 1)
	value := entries.NewValue([]byte("value1"))

	result := builder.Add(key, value)
	assert.True(t, result, "Expected Add to return true")
	assert.Equal(t, 1, len(builder.keyValueBeginOffset), "Expected one key-value pair to be added")
	assert.Equal(t, key, builder.firstKey, "Expected firstKey to be set to the added key")
}

func TestBuilder_Add_BlockSizeLimit(t *testing.T) {
	builder := NewBuilder(10) // Small block size to trigger limit
	key := entries.NewKey([]byte("key1"), 1)
	value := entries.NewValue([]byte("value1"))

	result := builder.Add(key, value)
	assert.False(t, result, "Expected Add to return false due to block size limit")
	assert.Equal(t, 0, len(builder.keyValueBeginOffset), "Expected no key-value pair to be added")
}

func TestBuilder_Add_FirstKeyAssignment(t *testing.T) {
	builder := NewBuilder(DefaultBlockSize)
	key1 := entries.NewKey([]byte("key1"), 1)
	value1 := entries.NewValue([]byte("value1"))
	key2 := entries.NewKey([]byte("key2"), 2)
	value2 := entries.NewValue([]byte("value2"))

	builder.Add(key1, value1)
	builder.Add(key2, value2)

	assert.Equal(t, key1, builder.firstKey, "Expected firstKey to be set to the first added key")
}

func TestBuilder_Add_MultipleKeys(t *testing.T) {
	builder := NewBuilder(DefaultBlockSize)
	key1 := entries.NewKey([]byte("key1"), 1)
	value1 := entries.NewValue([]byte("value1"))
	key2 := entries.NewKey([]byte("key2"), 2)
	value2 := entries.NewValue([]byte("value2"))

	result1 := builder.Add(key1, value1)
	result2 := builder.Add(key2, value2)

	assert.True(t, result1, "Expected first Add to return true")
	assert.True(t, result2, "Expected second Add to return true")
	assert.Equal(t, 2, len(builder.keyValueBeginOffset), "Expected two key-value pairs to be added")
}
