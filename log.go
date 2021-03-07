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
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
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

var buffers = sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 1024)) }}

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

	const messageKey, timestampKey = "message", "timestamp"
	var lenEncFields int
	for _, field := range encodedFields {
		key := field.Key()
		if key == messageKey || key == timestampKey {
			continue
		}
		lenEncFields += 2 + len(key) + 2 + len(field.Value())
	}

	const begin, msgInsert = `{ "` + timestampKey + `": "`, `", "` + messageKey + `": `
	sb := buffers.Get().(*bytes.Buffer)
	sb.Reset()
	sb.Grow(len(begin) + len(time.RFC3339) + len(msgInsert) + len(msg) + lenEncFields + 2)
	sb.WriteString(begin)
	sb.WriteString(now.Format(time.RFC3339))
	sb.WriteString(msgInsert)
	_ = json.NewEncoder(sb).Encode(msg)
	if sb.Bytes()[sb.Len()-1] == '\n' {
		sb.Truncate(sb.Len() - 1)
	}

	for _, field := range encodedFields {
		key := field.Key()
		if key == "message" || key == "timestamp" {
			continue
		}
		sb.WriteString(", ")
		sb.WriteString(key)
		sb.WriteString(`: `)
		sb.WriteString(field.Value())
	}
	sb.WriteString(" }\n")

	w := a.Writer
	if w == nil {
		w = os.Stderr
	}
	_, _ = w.Write(sb.Bytes())
	buffers.Put(sb)
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
