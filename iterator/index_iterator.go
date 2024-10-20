package iterator

import (
	"container/heap"
	"lsm/entries"
)

type MinHeapIndexIterator []IndexIterator

func (h *MinHeapIndexIterator) Len() int {
	return len(*h)
}

func (h *MinHeapIndexIterator) Less(i, j int) bool {
	return (*h)[i].Compare((*h)[j])
}

func (h *MinHeapIndexIterator) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeapIndexIterator) Push(x any) {
	*h = append(*h, x.(IndexIterator))
}

func (h *MinHeapIndexIterator) Pop() any {
	old := *h
	n := len(old)
	last := old[n-1]
	*h = old[0 : n-1]
	return last
}

type IndexIterator struct {
	index int
	Iterator
}

func NewIndexIterator(index int, it Iterator) *IndexIterator {
	return &IndexIterator{
		index:    index,
		Iterator: it,
	}
}

func (i IndexIterator) Compare(j IndexIterator) bool {
	comparison := i.Iterator.Key().Compare(j.Iterator.Key())
	if comparison == 0 {
		return i.index < j.index
	}
	return comparison < 0
}

type MergeIterator struct {
	current   IndexIterator
	iterators *MinHeapIndexIterator
}

func NewMergeIterator(iterators []Iterator) *MergeIterator {
	minHeap := &MinHeapIndexIterator{}
	heap.Init(minHeap)

	for i, it := range iterators {
		if it.IsValid() {
			heap.Push(minHeap, *NewIndexIterator(i, it))
		}
	}

	return &MergeIterator{
		current:   heap.Pop(minHeap).(IndexIterator),
		iterators: minHeap,
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
	for _, iter := range *it.iterators {
		if curr.Key().Compare(iter.Key()) != 0 {
			break
		}
		if !iter.Next() || !iter.IsValid() {
			heap.Pop(it.iterators)
			return false
		}
	}

	if !curr.Next() || !curr.IsValid() {
		if it.iterators.Len() > 0 {
			it.current = heap.Pop(it.iterators).(IndexIterator)
			return true
		}
		return false
	}

	if it.iterators.Len() > 0 {
		iter := *it.iterators
		if iter[0].Key().Compare(curr.Key()) < 0 {
			curr, iter[0] = iter[0], curr
			it.current = curr
			it.iterators = &iter
		}
	}
	return true
}
