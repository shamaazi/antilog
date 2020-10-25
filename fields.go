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
	var token struct{}
	// map for uniqueness
	m := make(map[string]struct{}, len(f)+len(fields))
	for _, v := range f {
		m[v.Key()] = token
	}
	var rev EncodedFields
	for ix := len(fields) - 1; ix >= 0; ix-- {
		field := fields[ix]
		key := field.Key()
		if _, found := m[key]; found {
			continue
		}
		// rev contains the reversed EncodedFields, to allow appends.
		if rev == nil {
			length := len(f)
			rev = make(EncodedFields, length, length+len(fields))
			for i, v := range f {
				rev[length-1-i] = v
			}
		}
		rev = append(rev, field)
		m[key] = token
	}
	if rev == nil {
		// Nothing new has been added
		return f
	}
	// Reverse the reversed array
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}
