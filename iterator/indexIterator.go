package iterator

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

func NewIndexIterator(index int, it Iterator) IndexIterator {
	return IndexIterator{
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
