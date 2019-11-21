package antilog_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestOmitsUnknownTypes(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test", "antilog", logger)
	logLine := parseLogLine(buffer.Bytes())

	expected := map[string]interface{}{
		"Fields": nil,
		"Writer": map[string]interface{}{},
	}
	require.EqualValues(t, expected, logLine["antilog"])
}

func TestIncludesContextFields(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer).With("test", "hello")

	logger.Write("this is a test")
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, "hello", logLine["test"])
}

func TestAppendsLoggedFieldsToContextFields(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer).With("test", "hello")

	logger.Write("this is a test", "tomato", "banana")
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, "banana", logLine["tomato"])
}

func TestPicksLastDuplicateValue(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test", "tomato", 1, "potato", 2, "pineapple", 3, "potato", 4)
	require.NotContains(t, buffer.String(), `"potato": 2`)
	require.Contains(t, buffer.String(), `"potato": 4`)

	logLine := parseLogLine(buffer.Bytes())
	require.EqualValues(t, 4, logLine["potato"])
}

func TestOverridesContextValue(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer).With("potato", 2)

	logger.Write("this is a test", "tomato", 1, "pineapple", 3, "potato", 4)
	require.NotContains(t, buffer.String(), `"potato": 2`)
	require.Contains(t, buffer.String(), `"potato": 4`)

	logLine := parseLogLine(buffer.Bytes())
	require.EqualValues(t, 4, logLine["potato"])
}

func TestReplacesContextValue(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer).With("potato", 2)

	logger = logger.With("potato", 4)

	logger.Write("this is a test", "tomato", 1, "pineapple", 3)
	require.NotContains(t, buffer.String(), `"potato": 2`)
	require.Contains(t, buffer.String(), `"potato": 4`)

	logLine := parseLogLine(buffer.Bytes())
	require.EqualValues(t, 4, logLine["potato"])
}

func TestLogsErrors(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	logger.Write("this is a test", "error", errors.New("an error occurred"))
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, "an error occurred", logLine["error"])
}

func TestLogsNilErrors(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	var err error
	logger.Write("this is a test", "error", err)
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, nil, logLine["error"])
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

func TestAlteringMapsDoesNotChangeLog(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer)

	values := map[string]string{
		"woo": "yay",
	}
	logger = logger.With("values", values)

	values["woo"] = "no"
	values["yay"] = "yes"

	logger.Write("this is a test")
	logLine := parseLogLine(buffer.Bytes())

	require.EqualValues(t, map[string]interface{}{"woo": "yay"}, logLine["values"])
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

func BenchmarkLogWithComplexFieldsInContext(b *testing.B) {
	buffer := &bytes.Buffer{}
	logger := antilog.WithWriter(buffer).With("complex field", map[string]interface{}{
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
	})

	for n := 0; n < b.N; n++ {
		logger.Write("a message", "simple field", "test")
	}
}
