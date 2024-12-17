package entries

type Kind int

const (
	KindPut Kind = iota
	KindDelete
)

type Entry struct {
	Key
	Value
	Kind
}

func (e Entry) IsKindPut() bool {
	return e.Kind == KindPut
}

func (e Entry) IsKindDelete() bool {
	return e.Kind == KindDelete
}

func (e Entry) SizeInBytes() int {
	return e.Key.EncodedKeySizeInBytes() + e.Value.Size()
}

type TimestampedBatch struct {
	entries []Entry
}

func NewTimestampedBatch(batch Batch, commitTimestamp uint64) *TimestampedBatch {
	timestampedBatch := TimestampedBatch{}
	for _, pair := range batch.pairs {
		if pair.kind == KindPut {
			timestampedBatch.put(NewKey(pair.key, commitTimestamp), pair.value)
		} else {
			timestampedBatch.delete(NewKey(pair.key, commitTimestamp))
		}
	}
	return &timestampedBatch
}

func (b *TimestampedBatch) AllEntries() []Entry {
	return b.entries
}

func (b *TimestampedBatch) SizeInBytes() int {
	size := 0
	for _, entry := range b.entries {
		size += entry.SizeInBytes()
	}
	return size
}

func (b *TimestampedBatch) put(key Key, value Value) *TimestampedBatch {
	b.entries = append(b.entries, Entry{key, value, KindPut})
	return b
}

func (b *TimestampedBatch) delete(key Key) *TimestampedBatch {
	b.entries = append(b.entries, Entry{key, EmptyValue, KindDelete})
	return b
}
