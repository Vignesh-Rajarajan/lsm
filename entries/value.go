package entries

type Value struct {
	Value []byte
}

var EmptyValue = Value{}

func NewValue(value []byte) Value {
	return Value{Value: value}
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

func (val Value) Size() int {
	return len(val.Value)
}

func (val Value) EncodeTo(buffer []byte) int {
	copy(buffer, val.Value)
	return len(val.Value)
}

func (val Value) DecodeFrom(buffer []byte) int {
	val.Value = buffer
	return len(val.Value)
}

func (val Value) Bytes() []byte {
	return val.Value
}
