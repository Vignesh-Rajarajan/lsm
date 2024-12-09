package iterator

import (
	"github.com/stretchr/testify/assert"
	"lsm/entries"
	"testing"
)

type testIterator struct {
	keys   []entries.Key
	values []entries.Value
	index  int
}

func newTestIterator(keys []entries.Key, values []entries.Value) *testIterator {
	return &testIterator{
		keys:   keys,
		values: values,
		index:  0,
	}
}

func (it *testIterator) Key() entries.Key {
	return it.keys[it.index]
}

func (it *testIterator) Value() entries.Value {
	return it.values[it.index]
}

func (it *testIterator) Next() bool {
	it.index++
	return true
}

func (it *testIterator) IsValid() bool {
	return it.index < len(it.keys)
}

func TestMergeIterator_SingleIterator(t *testing.T) {
	iter := newTestIterator([]entries.Key{entries.NewStringKey("a"), entries.NewStringKey("b")}, []entries.Value{entries.NewStringValue("1"), entries.NewStringValue("2")})
	merge := NewMergeIterator([]Iterator{iter})
	assert.True(t, merge.IsValid())
	assert.Equal(t, "a", string(merge.Key().Key))
	assert.Equal(t, "1", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "b", string(merge.Key().Key))
	assert.Equal(t, "2", string(merge.Value().Value))
	merge.Next()
	assert.False(t, merge.IsValid())
}

func TestMergeIterator_MultipleIterator(t *testing.T) {
	iter1 := newTestIterator([]entries.Key{entries.NewStringKey("a"), entries.NewStringKey("c")}, []entries.Value{entries.NewStringValue("1"), entries.NewStringValue("3")})
	iter2 := newTestIterator([]entries.Key{entries.NewStringKey("b"), entries.NewStringKey("d")}, []entries.Value{entries.NewStringValue("2"), entries.NewStringValue("4")})
	merge := NewMergeIterator([]Iterator{iter1, iter2})
	assert.True(t, merge.IsValid())
	assert.Equal(t, "a", string(merge.Key().Key))
	assert.Equal(t, "1", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "b", string(merge.Key().Key))
	assert.Equal(t, "2", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "c", string(merge.Key().Key))
	assert.Equal(t, "3", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "d", string(merge.Key().Key))
	assert.Equal(t, "4", string(merge.Value().Value))
	merge.Next()
	assert.False(t, merge.IsValid())
}

func TestMergeIterator_MultipleIteratorHavingSameKey(t *testing.T) {
	iter1 := newTestIterator([]entries.Key{entries.NewStringKey("key"), entries.NewStringKey("key1"), entries.NewStringKey("key2")}, []entries.Value{entries.NewStringValue("1"), entries.NewStringValue("2"), entries.NewStringValue("3")})
	iter2 := newTestIterator([]entries.Key{entries.NewStringKey("key"), entries.NewStringKey("test")}, []entries.Value{entries.NewStringValue("2"), entries.NewStringValue("test")})
	merge := NewMergeIterator([]Iterator{iter1, iter2})
	assert.True(t, merge.IsValid())
	assert.Equal(t, "key", string(merge.Key().Key))
	assert.Equal(t, "1", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "key1", string(merge.Key().Key))
	assert.Equal(t, "2", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "key2", string(merge.Key().Key))
	assert.Equal(t, "3", string(merge.Value().Value))
	merge.Next()
	assert.True(t, merge.IsValid())
	assert.Equal(t, "test", string(merge.Key().Key))
	assert.Equal(t, "test", string(merge.Value().Value))
}

func TestIndexIterator(t *testing.T) {
	idIteratorOne := NewIndexIterator(0, newTestIterator(
		[]entries.Key{entries.NewStringKey("a")},
		[]entries.Value{entries.NewStringValue("1")},
	))

	idIteratorTwo := NewIndexIterator(0, newTestIterator(
		[]entries.Key{entries.NewStringKey("b")},
		[]entries.Value{entries.NewStringValue("2")},
	))

	assert.True(t, idIteratorOne.Compare(idIteratorTwo))

	idIteratorOne = NewIndexIterator(0, newTestIterator(
		[]entries.Key{entries.NewStringKey("a")},
		[]entries.Value{entries.NewStringValue("1")},
	))

	idIteratorTwo = NewIndexIterator(0, newTestIterator(
		[]entries.Key{entries.NewStringKey("a")},
		[]entries.Value{entries.NewStringValue("2")},
	))

	assert.False(t, idIteratorOne.Compare(idIteratorTwo))
}

func TestMergeIterator_InValid(t *testing.T) {
	iterator :=
		newTestIterator([]entries.Key{entries.NewStringKey("a"),
			entries.NewStringKey("b")},
			[]entries.Value{
				entries.NewStringValue("1"),
				entries.NewStringValue("2"),
			})
	iterator.index = 2
	merge := NewMergeIterator([]Iterator{iterator})
	assert.False(t, merge.IsValid())

}
