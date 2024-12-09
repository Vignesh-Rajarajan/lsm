package entries

type Value struct {
	Value []byte
}

var EmptyValue = Value{}

func NewValue(value []byte) *Value {
	return &Value{Value: value}
}

func (val Value) String() string {
	return string(val.Value)
}

func NewStringValue(value string) Value {
	return Value{Value: []byte(value)}
}

func (val Value) IsEmpty() bool {
	return len(val.Value) == 0
}
