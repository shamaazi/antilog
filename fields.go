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
	if len(fields) == 0 {
		return f
	}
	var offset int
	for ix := len(fields) - 1; ix >= 0; ix-- {
		field := fields[ix]
		key := field.Key()
		var found bool
		for _, v := range f[offset:] {
			if v.Key() == key {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if offset == 0 {
			offset = len(fields)
			f = append(make(EncodedFields, offset, offset+len(f)), f...)
		}
		offset--
		f[offset] = field
	}
	return f[offset:]
}
