package entries

type Value struct {
	Value []byte
}

func (val Value) IsEmpty() bool {
	return val.Value == nil
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
