package antilog

// Field type for all inputs
type Field interface{}

// EncodedField type for storing fields in after conversion to JSON
type EncodedField [2]string

// Key of the encoded field
func (f EncodedField) Key() string {
	return f[0]
}

// Value of the encoded field
func (f EncodedField) Value() string {
	return f[1]
}

// EncodedFields is a list of encoded fields
type EncodedFields []EncodedField

// PrependUnique encoded field if the key is not already set
func (f EncodedFields) PrependUnique(fields EncodedFields) EncodedFields {
	for ix := len(fields) - 1; ix >= 0; ix-- {
		field := fields[ix]

		var found bool
		for _, v := range f {
			if v.Key() == field.Key() {
				found = true
				break
			}
		}

		if !found {
			f = append(EncodedFields{field}, f...)
		}
	}

	return f
}
