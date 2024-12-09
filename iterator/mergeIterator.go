package iterator

import (
	"container/heap"
	"lsm/entries"
)

type MergeIterator struct {
	current   IndexIterator
	iterators *MinHeapIndexIterator
}

func NewMergeIterator(iterators []Iterator) *MergeIterator {
	minHeap := &MinHeapIndexIterator{}
	heap.Init(minHeap)

	for i, it := range iterators {
		if it.IsValid() {
			heap.Push(minHeap, NewIndexIterator(i, it))
		}
	}

	if minHeap.Len() > 0 {
		return &MergeIterator{
			current:   heap.Pop(minHeap).(IndexIterator),
			iterators: minHeap,
		}
	}

	return &MergeIterator{
		current: NewIndexIterator(0, EmptyIterator),
	}
}

func (it *MergeIterator) Key() entries.Key {
	return it.current.Key()
}

func (it *MergeIterator) Value() entries.Value {
	return it.current.Value()
}

func (it *MergeIterator) IsValid() bool {
	return it.current.IsValid()
}

// Next advances the MergeIterator to the next valid key-value pair. It returns true if the
// iterator was successfully advanced, and false if there are no more valid key-value pairs.
// Next() will return false if any of the underlying iterators become invalid.
func (it *MergeIterator) Next() bool {
	curr := it.current
	if !it.advanceOtherIteratorForSameKey(curr) {
		return false
	}

	if !it.advance(curr) {
		return false
	}

	if it.maybePop(curr) {
		return true
	}

	return it.maybeSwapCurrent(curr)
}

func (it *MergeIterator) advanceOtherIteratorForSameKey(curr IndexIterator) bool {
	for _, iter := range *it.iterators {
		if curr.Key().EqualTo(iter.Key()) {
			if !it.advance(iter) {
				heap.Pop(it.iterators)
				return false
			}
			break
		}
	}
	return true
}

func (it *MergeIterator) maybePop(current IndexIterator) bool {
	if !current.IsValid() {
		if it.iterators.Len() > 0 {
			it.current = heap.Pop(it.iterators).(IndexIterator)
		}
		return true
	}
	return false
}

func (it *MergeIterator) maybeSwapCurrent(current IndexIterator) bool {
	if it.iterators.Len() > 0 {
		iter := *it.iterators
		if !current.Compare(iter[0]) {
			current, iter[0] = iter[0], current
			it.current = current
			it.iterators = &iter
		}
	}
	return true
}

func (it *MergeIterator) advance(idxIterator IndexIterator) bool {
	return idxIterator.Next()
}
