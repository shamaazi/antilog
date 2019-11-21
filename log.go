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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// AntiLog is the antidote to modern loggers
type AntiLog struct {
	Fields EncodedFields
	Writer io.Writer
}

// New instance of AntiLog
func New() AntiLog {
	return AntiLog{}
}

// With returns a copy of the AntiLog instance with the provided fields preset for every subsequent call.
func (a AntiLog) With(fields ...Field) AntiLog {
	a.Fields = encodeFieldList(fields).PrependUnique(a.Fields)
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

	encodedFields := EncodedFields{}.
		PrependUnique(encodeFieldList(fields)).
		PrependUnique(a.Fields)

	for ix := 0; ix < len(encodedFields); ix++ {
		field := encodedFields[ix]
		if field.Key() != "message" || field.Key() != "timestamp" {
			continue
		}
		encodedFields = append(encodedFields[0:ix], encodedFields[ix+1:len(encodedFields)]...)
	}

	encodedFields = append(encodeFieldList([]Field{"timestamp", now.Format(time.RFC3339), "message", msg}), encodedFields...)

	if a.Writer == nil {
		a.Writer = os.Stderr
	}

	fmt.Fprintln(a.Writer, encodedFieldsToJSONObject(encodedFields))
}

func toJSON(field Field) string {
	// In the case of errors, explicitly destructure them
	if err, ok := field.(error); ok {
		field = err.Error()
	}

	// For anything else, just let json.Marshal do it
	bytes, err := json.Marshal(field)
	if err != nil {
		return string(err.Error())
	}

	return string(bytes)
}

func encodeFieldList(fields []Field) EncodedFields {
	convertedFields := make(EncodedFields, 0, len(fields))

	numFields := len(fields) / 2
	for ix := 0; ix < numFields; ix++ {
		rawKey := fields[ix*2]
		rawValue := fields[ix*2+1]

		keyString, ok := rawKey.(string)
		if !ok {
			continue
		}

		key := toJSON(keyString)
		value := toJSON(rawValue)

		convertedFields = append(convertedFields, EncodedField{key, value})
	}
	return convertedFields
}

func encodedFieldsToJSONObject(fields []EncodedField) string {
	var sb stringBuilder
	sb.WriteString(`{ `)

	var comma bool
	for _, field := range fields {
		if comma {
			sb.WriteString(`, `)
		}
		sb.WriteStrings(field.Key(), `: `, field.Value())
		comma = true
	}

	sb.WriteString(` }`)
	return sb.String()
}
