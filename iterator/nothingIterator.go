package iterator

import "lsm/entries"

type NothingIterator struct {
}

var EmptyIterator = &NothingIterator{}

func (it NothingIterator) Next() bool {
	return false
}

func (it NothingIterator) Key() entries.Key {
	return *entries.NewKey(nil)
}

func (it NothingIterator) Value() entries.Value {
	return entries.EmptyValue
}

func (it NothingIterator) IsValid() bool {
	return false
}
