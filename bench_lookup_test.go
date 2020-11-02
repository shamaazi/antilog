// +build lookup

package antilog_test

import (
	"strconv"
	"testing"
)

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
