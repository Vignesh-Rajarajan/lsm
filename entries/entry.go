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
