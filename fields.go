package antilog

// Field type for all inputs
type Field interface{}

// EncodedField type for storing fields in after conversion to JSON
type EncodedField string

func (f EncodedField) String() string {
	return string(f)
}
