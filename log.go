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

// Field type for all inputs
type Field interface{}

// AntiLog is the antidote to modern loggers
type AntiLog struct {
	Fields []Field
	Writer io.Writer
}

// With returns a copy of the AntiLog instance with the provided fields preset for every subsequent call.
func (a AntiLog) With(fields ...Field) AntiLog {
	a.Fields = append(a.Fields, fields...)
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
	combinedFields := []Field{
		"timestamp", now.Format(time.RFC3339),
		"message", msg,
	}
	combinedFields = append(combinedFields, a.Fields...)
	combinedFields = append(combinedFields, fields...)

	if a.Writer == nil {
		a.Writer = os.Stderr
	}

	fmt.Fprintln(a.Writer, toJSONObject(combinedFields))
}
