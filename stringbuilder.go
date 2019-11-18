package antilog

import (
	"fmt"
	"strings"
)

type stringBuilder struct {
	strings.Builder
}

func (b *stringBuilder) WriteStrings(strings ...string) {
	for _, s := range strings {
		b.WriteString(s)
	}
}

func (b *stringBuilder) WriteStringf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format, args...))
}
