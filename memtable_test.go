package main

import (
	"github.com/stretchr/testify/assert"
	"lsm/entries"
	"lsm/txn"
	"testing"
)

func TestMemtable(t *testing.T) {
	memtable := NewMemTable(1)
	assert.Equal(t, uint64(0), memtable.size.Load())

	key := entries.NewStringKey("key")
	value := entries.NewStringValue("value")
	memtable.Set(key, value)

	assert.NotEmpty(t, memtable.entries)
	value, ok := memtable.Get(key)
	assert.True(t, ok)
	assert.Equal(t, "value", string(value.Value))
	assert.Equal(t, uint64(8), memtable.size.Load())
	val, ok := memtable.Get(entries.NewStringKey("test"))
	assert.False(t, ok)
	assert.Equal(t, entries.EmptyValue, val)

	memtable.Delete(key)
	value, ok = memtable.Get(key)
	assert.False(t, ok)
	assert.Equal(t, entries.EmptyValue, value)
}

func TestMemtable_ScanInclusive(t *testing.T) {
	memtable := NewMemTable(1)
	memtable.Set(entries.NewStringKey("key1"), entries.NewStringValue("value1"))
	memtable.Set(entries.NewStringKey("key2"), entries.NewStringValue("value2"))
	memtable.Set(entries.NewStringKey("key3"), entries.NewStringValue("value3"))

	iterator := memtable.Scan(txn.NewInclusiveRange(entries.NewStringKey("key2"), entries.NewStringKey("key2")))
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "value2", string(iterator.Value().Value))
	iterator.Next()
	assert.False(t, iterator.IsValid())

	iterator = memtable.Scan(txn.NewInclusiveRange(entries.NewStringKey("key2"), entries.NewStringKey("key6")))
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "value2", string(iterator.Value().Value))
	iterator.Next()
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "value3", string(iterator.Value().Value))
	iterator.Next()
	assert.False(t, iterator.IsValid())
}
