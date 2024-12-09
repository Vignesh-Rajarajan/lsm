package main

import (
	"lsm/entries"
	"lsm/iterator"
	"lsm/txn"
)

type StorageOptions struct {
	memTableSize uint64
}

type StorageState struct {
	memtable           *MemTable
	immutableMemtables []*MemTable
	options            StorageOptions
}

func NewStorageState() *StorageState {
	return &StorageState{
		memtable: NewMemTable(1),
		options:  StorageOptions{memTableSize: 1 << 20},
	}
}

func NewStorageStateWithOptions(options StorageOptions) *StorageState {
	return &StorageState{
		memtable: NewMemTable(1),
		options:  options,
	}
}

func (s *StorageState) Get(key entries.Key) (entries.Value, bool) {
	val, ok := s.memtable.Get(key)
	if ok {
		return val, ok
	}

	for i := len(s.immutableMemtables) - 1; i >= 0; i-- {
		val, ok = s.immutableMemtables[i].Get(key)
		if ok {
			return val, ok
		}
	}
	return entries.EmptyValue, false
}

func (s *StorageState) Set(batch *entries.Batch) {
	for _, entry := range batch.AllEntries() {
		if entry.IsKindPut() {
			s.memtable.Set(entry.Key, entry.Value)
		} else {
			s.memtable.Delete(entry.Key)
		}
	}

	s.couldFreezeMemtable()
}

func (s *StorageState) couldFreezeMemtable() {
	if s.memtable.Size() >= s.options.memTableSize {
		s.immutableMemtables = append(s.immutableMemtables, s.memtable)
		s.memtable = NewMemTable(1)
	}
}

func (s *StorageState) forceFreezeMemtable() {
	if !s.memtable.IsEmpty() {
		s.immutableMemtables = append(s.immutableMemtables, s.memtable)
		s.memtable = NewMemTable(1)
	}
}

func (s *StorageState) HasImmutableTables() bool {
	return len(s.immutableMemtables) > 0
}

func (s *StorageState) Scan(inclusiveRange txn.InclusiveRange) iterator.Iterator {

	iterators := make([]iterator.Iterator, len(s.immutableMemtables)+1)
	iterators[0] = s.memtable.Scan(inclusiveRange)
	index := 1

	for immutableMemIndex := len(s.immutableMemtables) - 1; immutableMemIndex >= 0; immutableMemIndex-- {
		iterators[index] = s.immutableMemtables[immutableMemIndex].Scan(inclusiveRange)
		index++
	}

	return iterator.NewMergeIterator(iterators)
}
