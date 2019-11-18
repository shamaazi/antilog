package antilog

import (
	"fmt"
	"reflect"
	"strconv"
)

func extractArrayAsValues(v reflect.Value) []Field {
	values := []Field{}

	for ix := 0; ix < v.Len(); ix++ {
		values = append(values, v.Index(ix).Interface())
	}

	return values
}

func extractStructAsFields(v reflect.Value) []Field {
	fields := []Field{}
	t := v.Type()
	for ix := 0; ix < v.NumField(); ix++ {
		fields = append(fields, t.Field(ix).Name, v.Field(ix).Interface())
	}
	return fields
}

func extractMapAsFields(v reflect.Value) []Field {
	fields := []Field{}
	iter := v.MapRange()
	for iter.Next() {
		subkey := iter.Key()
		subvalue := iter.Value()

		if key, ok := subkey.Interface().(string); ok {
			fields = append(fields, key, subvalue.Interface())
		}
	}
	return fields
}

func toJSONObjectFields(fields []Field) string {
	var sb stringBuilder

	numFields := len(fields) / 2
	var comma bool
	for ix := 0; ix < numFields; ix++ {
		rawKey := fields[ix*2]
		rawValue := fields[ix*2+1]

		var key string
		var ok bool
		if key, ok = rawKey.(string); !ok {
			continue
		}

		if value, ok := toJSON(rawValue); ok {
			if comma {
				sb.WriteString(`, `)
			}

			sb.WriteStrings(strconv.Quote(key), `: `, value.String())
			comma = true
		}
	}

	return sb.String()
}

func toJSONObject(fields []Field) EncodedField {
	var sb stringBuilder
	sb.WriteString(`{ `)

	sb.WriteString(toJSONObjectFields(fields))

	sb.WriteString(` }`)
	return EncodedField(sb.String())
}

func toJSONArray(values []Field) EncodedField {
	var sb stringBuilder
	sb.WriteString(`[ `)

	var comma bool
	for ix := 0; ix < len(values); ix++ {
		if value, ok := toJSON(values[ix]); ok {
			if comma {
				sb.WriteString(`, `)
			}
			sb.WriteString(value.String())
			comma = true
		}
	}

	sb.WriteString(` ]`)
	return EncodedField(sb.String())
}

func toJSON(field Field) (EncodedField, bool) {
	v := reflect.ValueOf(field)

	var value EncodedField
	switch v.Kind() {
	case reflect.Bool:
		value = EncodedField(fmt.Sprintf("%v", v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		value = EncodedField(fmt.Sprintf("%v", v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		value = EncodedField(fmt.Sprintf("%v", v.Uint()))
	case reflect.Float32, reflect.Float64:
		value = EncodedField(fmt.Sprintf("%v", v.Float()))
	case reflect.String:
		value = EncodedField(strconv.Quote(v.String()))
	case reflect.Slice:
		values := extractArrayAsValues(v)
		value = toJSONArray(values)
	case reflect.Map:
		subfields := extractMapAsFields(v)
		value = toJSONObject(subfields)
	case reflect.Struct:
		subfields := extractStructAsFields(v)
		value = toJSONObject(subfields)
	default:
		return "", false
	}

	return value, true
}
