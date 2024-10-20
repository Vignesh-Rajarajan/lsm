package main

import (
	"bytes"
	"github.com/huandu/skiplist"
	"lsm/entries"
	"lsm/iterator"
	"sync/atomic"
)

type MemTable struct {
	id      uint
	size    atomic.Uint64
	entries *skiplist.SkipList
}

func NewMemTable(id uint) *MemTable {
	return &MemTable{
		id: id,
		entries: skiplist.New(skiplist.GreaterThanFunc(func(key, otherKey interface{}) int {
			left := key.(entries.Key)
			right := otherKey.(entries.Key)
			return bytes.Compare(left.Key, right.Key)
		})),
	}
}

func (m *MemTable) Set(key entries.Key, value entries.Value) {
	size := uint64(len(key.Key) + len(value.Value))
	m.size.Add(size)
	m.entries.Set(key, value)
}

func (m *MemTable) Get(key entries.Key) (entries.Value, bool) {
	value, ok := m.entries.GetValue(key)
	if !ok || value.(entries.Value).IsEmpty() {
		return entries.EmptyValue, false
	}
	return value.(entries.Value), true
}

func (m *MemTable) IsEmpty() bool {
	return m.entries.Len() == 0
}

func (m *MemTable) Size() uint64 {
	return m.size.Load()
}

func (m *MemTable) Delete(key entries.Key) {
	m.Set(key, entries.EmptyValue)
}

func (m *MemTable) ScanInclusive(start, end entries.Key) *iterator.MemTableIterator {
	return iterator.NewMemtableIterator(m.entries.Find(start), end)
}
