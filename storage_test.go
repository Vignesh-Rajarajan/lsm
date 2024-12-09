package main

import (
	"github.com/stretchr/testify/assert"
	"lsm/entries"
	"lsm/txn"
	"testing"
)

func TestStorageState_HasImmutableTables(t *testing.T) {
	state := NewStorageState()
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("key"), entries.NewStringValue("value"))
	state.Set(batch)
	assert.False(t, state.HasImmutableTables())
}

func TestStorageState_Get(t *testing.T) {
	state := NewStorageState()
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("key"), entries.NewStringValue("value"))
	batch.Put(entries.NewStringKey("key2"), entries.NewStringValue("value2"))
	batch.Put(entries.NewStringKey("key3"), entries.NewStringValue("value3"))

	state.Set(batch)

	val, ok := state.Get(entries.NewStringKey("key"))
	assert.True(t, ok)
	assert.Equal(t, "value", string(val.Value))
	val, ok = state.Get(entries.NewStringKey("key2"))
	assert.True(t, ok)
	assert.Equal(t, "value2", string(val.Value))
	val, ok = state.Get(entries.NewStringKey("key3"))
	assert.True(t, ok)
	assert.Equal(t, "value3", string(val.Value))
}

func TestStorage_Delete(t *testing.T) {
	state := NewStorageState()
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("key"), entries.NewStringValue("value"))
	state.Set(batch)
	batch = entries.NewBatch()
	batch.Delete(entries.NewStringKey("key"))
	state.Set(batch)

	value, ok := state.Get(entries.NewStringKey("key"))
	assert.False(t, ok)
	assert.Equal(t, entries.EmptyValue, value)
}

func TestNewStorageStateWithOptions(t *testing.T) {
	options := StorageOptions{memTableSize: 10}
	state := NewStorageStateWithOptions(options)
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("SOCK_STREAM"), entries.NewStringValue("Stream Sockets"))
	batch.Put(entries.NewStringKey("SOCK_DGRAM"), entries.NewStringValue("Datagram sockets"))
	batch.Put(entries.NewStringKey("key3"), entries.NewStringValue("value3"))
	state.Set(batch)
	assert.True(t, state.HasImmutableTables())

	value, ok := state.Get(entries.NewStringKey("SOCK_STREAM"))
	assert.True(t, ok)
	assert.Equal(t, "Stream Sockets", string(value.Value))
}

func TestStorageState_Scan(t *testing.T) {
	state := NewStorageState()
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("consensus"), entries.NewStringValue("value1"))
	batch.Put(entries.NewStringKey("storage"), entries.NewStringValue("value2"))
	batch.Put(entries.NewStringKey("data-structure"), entries.NewStringValue("value3"))
	state.Set(batch)

	iterator := state.Scan(txn.NewInclusiveRange(entries.NewStringKey("accurate"), entries.NewStringKey("etcd")))
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "consensus", string(iterator.Key().Key))
	assert.Equal(t, "value1", string(iterator.Value().Value))
	iterator.Next()
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "data-structure", string(iterator.Key().Key))
	assert.Equal(t, "value3", string(iterator.Value().Value))
	iterator.Next()
	assert.False(t, iterator.IsValid())
}

func TestStorageState_Scan_Inclusive(t *testing.T) {
	state := NewStorageState()
	batch := entries.NewBatch()
	batch.Put(entries.NewStringKey("consensus"), entries.NewStringValue("value1"))
	batch.Put(entries.NewStringKey("storage"), entries.NewStringValue("value2"))
	batch.Put(entries.NewStringKey("data-structure"), entries.NewStringValue("value3"))
	state.Set(batch)

	iterator := state.Scan(txn.NewInclusiveRange(entries.NewStringKey("accurate"), entries.NewStringKey("data-structure")))
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "consensus", string(iterator.Key().Key))
	assert.Equal(t, "value1", string(iterator.Value().Value))
	iterator.Next()
	assert.True(t, iterator.IsValid())
	assert.Equal(t, "data-structure", string(iterator.Key().Key))
	assert.Equal(t, "value3", string(iterator.Value().Value))
	iterator.Next()
	assert.False(t, iterator.IsValid())
}
