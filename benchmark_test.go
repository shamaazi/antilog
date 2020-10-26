package antilog

import (
	"io/ioutil"
	"testing"
	"time"
)

var (
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogEmpty(b *testing.B) {
	logger := New()
	logger.Writer = ioutil.Discard
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Write("")
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	logger := WithWriter(ioutil.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Write(fakeMessage)
		}
	})
}

func BenchmarkContextFields(b *testing.B) {
	logger := WithWriter(ioutil.Discard).With(
		"string", "four!",
		"time", time.Time{},
		"int", 123,
		"float", -2.203230293249593)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Write(fakeMessage)
		}
	})
}

func BenchmarkContextAppend(b *testing.B) {
	logger := WithWriter(ioutil.Discard).With("foo", "bar")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Write("", "bar", "baz")
		}
	})
}

func BenchmarkLogFields(b *testing.B) {
	logger := WithWriter(ioutil.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Write(fakeMessage,
				"string", "four!",
				"time", time.Time{},
				"int", 123,
				"float", -2.203230293249593)
		}
	})
}
