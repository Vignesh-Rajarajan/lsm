package entries

type Batch struct {
	entries []Entry
}

func NewBatch() *Batch {
	return &Batch{}
}

func (b *Batch) Put(key Key, value Value) {
	b.entries = append(b.entries, Entry{key, value, KindPut})
}

func (b *Batch) Delete(key Key) {
	b.entries = append(b.entries, Entry{key, EmptyValue, KindDelete})
}

func (b *Batch) AllEntries() []Entry {
	return b.entries
}
