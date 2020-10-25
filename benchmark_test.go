package antilog

import (
	"io/ioutil"
	"strconv"
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

var lookupArray = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func BenchmarkLookupMap(b *testing.B) {
	for length := 1; length < len(lookupArray); length++ {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			b.StopTimer()
			slice := lookupArray[:length]
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				var token struct{}
				m := make(map[string]struct{}, len(slice))
				for _, s := range slice {
					m[s] = token
				}
				var n int
				for _, s := range lookupArray[:] {
					if _, found := m[s]; found {
						n++
					}
				}
				b.StopTimer()
			}
		})
	}
}
func BenchmarkLookupWalk(b *testing.B) {
	for length := 1; length < len(lookupArray); length++ {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			b.StopTimer()
			slice := lookupArray[:length]
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				var n int
				for _, s := range lookupArray[:] {
					var found bool
					for _, v := range slice {
						if v == s {
							found = true
							break
						}
					}
					if found {
						n++
					}
				}
				b.StopTimer()
			}
		})
	}
}
