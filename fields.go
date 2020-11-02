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
	var res EncodedFields
	offset := len(fields)
	for ix := len(fields) - 1; ix >= 0; ix-- {
		field := fields[ix]
		flds := f
		if res != nil {
			flds = res[offset:]
		}
		key := field.Key()
		var found bool
		for _, v := range flds {
			if v.Key() == key {
				found = true
				break
			}
		}
		if found {
			continue
		}
		// res contains the EncodedFields, starting at offset
		if res == nil {
			length := len(f)
			res = make(EncodedFields, offset+length)
			copy(res[offset:], f)
		}
		offset--
		res[offset] = field
	}
	if res == nil {
		// Nothing new has been added
		return f
	}
	return res[offset:]
}
