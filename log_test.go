package antilog_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/shamaazi/antilog"
	"github.com/stretchr/testify/require"
)

func parseLogLine(b []byte) (v map[string]interface{}) {
	err := json.Unmarshal(b, &v)
	if err != nil {
		panic(err)
	}
	return
}

func parseTime(i interface{}) (t time.Time) {
	s := i.(string)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return
}

func TestHasTimestampAndMessage(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test")

	logLine := parseLogLine(buffer.Bytes())

	require.Len(t, logLine, 2)
	require.WithinDuration(t, time.Now(), parseTime(logLine["timestamp"]), 1*time.Second)
	require.Equal(t, "this is a test", logLine["message"])
}

func TestHandlesBasicTypes(t *testing.T) {
	for test, value := range map[string]interface{}{
		"bool":   true,
		"int":    1234,
		"float":  123.456,
		"string": "wibble",
		"array":  []interface{}{"woo", "yay", "houpla"},
		"map":    map[string]interface{}{"woo": "yay", "houpla": "panowie"},
	} {
		t.Run(test, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			logger := antilog.WithWriter(buffer)

			logger.Write("this is a test", test, value)
			logLine := parseLogLine(buffer.Bytes())

			require.EqualValues(t, value, logLine[test])
		})
	}
}

type OuterStruct struct {
	Name         string
	Inner        InnerStruct
	AnotherField string
}

type InnerStruct struct {
	Tag          string
	ArrayOfStuff []LeafStruct
}

type LeafStruct struct {
	Key   string
	Value string
}

func TestHandlesNestedStructs(t *testing.T) {
	inputStructure := OuterStruct{
		Name: "Test struct",
		Inner: InnerStruct{
			Tag: "something",
			ArrayOfStuff: []LeafStruct{
				{"a key", "a value"},
				{"another", "with another value"},
				{"one more", "for luck"},
			},
		},
		AnotherField: "what is this?",
	}

	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test", "struct test", inputStructure)

	var actual struct {
		OuterStruct OuterStruct `json:"struct test"`
	}

	err := json.Unmarshal(buffer.Bytes(), &actual)
	require.NoError(t, err)

	require.EqualValues(t, inputStructure, actual.OuterStruct)
}

func TestHandlesDeeplyNestedTypes(t *testing.T) {
	inputStructure := map[string]interface{}{
		"array_with_various_types": []interface{}{
			"string",
			123.456,
			[]interface{}{
				"another",
				"array",
				"inside",
			},
			map[string]interface{}{
				"a map": "nested in the array",
			},
		},
		"map_with_various_types": map[string]interface{}{
			"string": "a string",
			"number": 1234.0,
			"bool":   false,
			"an array!": []interface{}{
				"with",
				"mixed",
				false,
				"types",
				map[string]interface{}{
					"including": "a map",
				},
			},
			"another map": map[string]interface{}{
				"with its own values": "like this",
			},
		},
	}

	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test", "a deep structure", inputStructure)
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, inputStructure, logLine["a deep structure"])
}

func BenchmarkLogWithNoFields(b *testing.B) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	for n := 0; n < b.N; n++ {
		logger.Write("a message")
	}
}

func BenchmarkLogWithComplexFields(b *testing.B) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)
	inputStructure := map[string]interface{}{
		"array_with_various_types": []interface{}{
			"string",
			123.456,
			[]interface{}{
				"another",
				"array",
				"inside",
			},
			map[string]interface{}{
				"a map": "nested in the array",
			},
		},
		"map_with_various_types": map[string]interface{}{
			"string": "a string",
			"number": 1234.0,
			"bool":   false,
			"an array!": []interface{}{
				"with",
				"mixed",
				false,
				"types",
				map[string]interface{}{
					"including": "a map",
				},
			},
			"another map": map[string]interface{}{
				"with its own values": "like this",
			},
		},
		"a struct of all things": struct {
			Name string
			Age  int
		}{"Mr Blobby", 48},
	}

	for n := 0; n < b.N; n++ {
		logger.Write("a message", "complex field", inputStructure)
	}
}
