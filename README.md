# antilog

Package antilog is the antidote to modern loggers.

* AntiLog only logs JSON formatted output. Structured logging is the only good logging.
* AntiLog does not have log levels. If you don't want something logged, don't log it.
* AntiLog does support setting fields in context. Useful for building a log context over the course of an operation.

## Basic Usage

```go
    antilog.Write("a message")
```

```json
{ "timestamp": "2019-11-18T14:00:32Z", "message": "a message" }`
```

## With Fields

```go
    antilog.Write("a message",
        "field", "value",
        "a_number", 123,
        "a_bool", false,
    )
```

```json
{ "timestamp": "2019-11-18T14:00:32Z", "message": "a message", "field": "value", "a_number": 123, "a_bool": false }`
```

## With Context

```go
    logger := antilog.With(
        "request_id", "12345",
        "user_id": "big_jim_mcdonald",
    )

    logger.Write("a message",
        "field", "value",
        "a_number", 123,
        "a_bool", false)
```

```json
{ "timestamp": "2019-11-18T14:00:32Z", "message": "a message", "request_id": "12345", "user_id": "big_jim_mcdonald", "field": "value", "a_number": 123, "a_bool": false }`
```

## Output Somewhere Else

```go
    var sb strings.Builder
    logger := antilog.WithWriter(sb)

    logger.Write("a message",
        "field", "value",
        "a_number", 123,
        "a_bool", false)

    fmt.Println(sb.String())
```

```json
{ "timestamp": "2019-11-18T14:00:32Z", "message": "a message", "field": "value", "a_number": 123, "a_bool": false }`
```
