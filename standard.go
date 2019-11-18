package antilog

import "io"

var antilog = AntiLog{}

// WithWriter returns a copy of the standard AntiLog instance configured to write to the given writer
func WithWriter(w io.Writer) AntiLog {
	return AntiLog{
		Writer: w,
	}
}

// With returns a copy of the standard AntiLog instance configured with the provided fields
func With(fields ...Field) AntiLog {
	return antilog.With(fields...)
}

// Write a message using the standard AntiLog instance
func Write(msg string, fields ...Field) {
	antilog.Write(msg, fields...)
}
