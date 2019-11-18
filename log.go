//Package antilog is the antidote to modern loggers.
//
// AntiLog only logs JSON formatted output. Structured logging is the only good
// logging.
//
// AntiLog does not have log levels. If you don't want something logged, don't
// log it.
//
// AntiLog does support setting fields in context. Useful for building a log
// context over the course of an operation.
package antilog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AntiLog is the antidote to modern loggers
type AntiLog struct {
	Fields []EncodedField
	Writer io.Writer
}

// With returns a copy of the AntiLog instance with the provided fields preset for every subsequent call.
func (a AntiLog) With(fields ...Field) AntiLog {
	a.Fields = append(a.Fields, encodeFieldList(fields)...)
	return a
}

// Write a JSON message to the configured writer or os.Stderr.
//
// Includes the message with the key `message`. Includes the timestamp with the
// key `timestamp`. The timestamp field is always first and the message second.
//
// Fields in context will not be overridden. AntiLog will log the same key
// multiple times if it is set multiple times. If you don't want that, don't
// specify it multiple times.
func (a AntiLog) Write(msg string, fields ...Field) {
	now := time.Now().UTC()
	combinedFields := encodeFieldList([]Field{
		"timestamp", now.Format(time.RFC3339),
		"message", msg,
	})
	combinedFields = append(combinedFields, a.Fields...)
	combinedFields = append(combinedFields, encodeFieldList(fields)...)

	if a.Writer == nil {
		a.Writer = os.Stderr
	}

	fmt.Fprintln(a.Writer, encodedFieldsToJSONObject(combinedFields))
}

func encodeFieldList(fields []Field) []EncodedField {
	convertedFields := make([]EncodedField, 0, len(fields))

	numFields := len(fields) / 2
	for ix := 0; ix < numFields; ix++ {
		rawKey := fields[ix*2]
		rawValue := fields[ix*2+1]

		keyString, ok := rawKey.(string)
		if !ok {
			continue
		}

		key, ok := toJSON(keyString)
		if !ok {
			continue
		}

		value, ok := toJSON(rawValue)
		if !ok {
			continue
		}

		convertedFields = append(convertedFields, key, value)
	}
	return convertedFields
}

func encodedFieldsToJSONObject(fields []EncodedField) string {
	var sb stringBuilder
	sb.WriteString(`{ `)

	numFields := len(fields) / 2
	var comma bool
	for ix := 0; ix < numFields; ix++ {
		key := fields[ix*2]
		value := fields[ix*2+1]

		if comma {
			sb.WriteString(`, `)
		}
		sb.WriteStrings(key.String(), `: `, value.String())
		comma = true
	}

	sb.WriteString(` }`)
	return sb.String()
}
