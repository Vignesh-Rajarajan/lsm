package iterator

import (
	"github.com/huandu/skiplist"
	"lsm/entries"
)

type Iterator interface {
	Next() bool
	Key() entries.Key
	Value() entries.Value
	IsValid() bool
}

type MemTableIterator struct {
	element *skiplist.Element
	endKey  entries.Key
}

func NewMemtableIterator(element *skiplist.Element, endKey entries.Key) *MemTableIterator {
	return &MemTableIterator{
		element: element,
		endKey:  endKey,
	}
}

func (it *MemTableIterator) Key() entries.Key {
	return it.element.Key().(entries.Key)
}

func (it *MemTableIterator) Value() entries.Value {
	return it.element.Value.(entries.Value)
}

func (it *MemTableIterator) Next() bool {
	el := it.element.Next()
	if el == nil {
		it.element = nil
		return false
	}
	key := el.Key().(entries.Key)
	if key.IsLessThanOrEqual(it.endKey) {
		it.element = el
		return true
	}
	it.element = nil
	return false
}

func (it *MemTableIterator) IsValid() bool {
	return it.element != nil
}
