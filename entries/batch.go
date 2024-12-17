package entries

import (
	"bytes"
	"errors"
)

type RawKeyValuePair struct {
	key   []byte
	value Value
	kind  Kind
}

func (kv RawKeyValuePair) Key() []byte {
	return kv.key
}

func (kv RawKeyValuePair) Value() Value {
	return kv.value
}

var DuplicateKeyErr = errors.New("batch already contains key")

type Batch struct {
	pairs []RawKeyValuePair
}

func NewBatch() *Batch {
	return &Batch{}
}

func (b *Batch) Put(key, value []byte) error {
	if b.Contains(key) {
		return DuplicateKeyErr
	}
	b.pairs = append(b.pairs, RawKeyValuePair{key, NewValue(value), KindPut})
	return nil

}

func (b *Batch) Delete(key []byte) {
	b.pairs = append(b.pairs, RawKeyValuePair{key: key, value: EmptyValue, kind: KindDelete})
}

func (b *Batch) Get(key []byte) (Value, bool) {
	for _, pair := range b.pairs {
		if bytes.Compare(pair.key, key) == 0 {
			return pair.value, true
		}
	}
	return EmptyValue, false
}

func (b *Batch) Contains(key []byte) bool {
	_, ok := b.Get(key)
	return ok
}

func (b *Batch) Len() int {
	return len(b.pairs)
}

func (b *Batch) CloneKeyValuePairs() []RawKeyValuePair {
	cloned := make([]RawKeyValuePair, 0, len(b.pairs))
	for _, pair := range b.pairs {
		cloned = append(cloned, pair)
	}
	return cloned
}
