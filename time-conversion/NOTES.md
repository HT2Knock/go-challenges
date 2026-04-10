# Time Conversion Challenge

## What Happens

Converts Unix timestamp `1585750374` to a TinyTime struct and adds 48 hours.

```text
1585750374 converted to a tinytime is: 2020-03-30 05:32:54 +0000 UTC
```

## What is Unix Timestamp?

Seconds elapsed since **January 1, 1970 00:00:00 UTC** (the Unix epoch).

- `1585750374` = **Sat Mar 28, 2020 05:32:54 UTC**
- Adding 48 hours = **Mon Mar 30, 2020 05:32:54 UTC**

## Why go-tinytime?

`time.Time` in Go is 24 bytes. TinyTime is ~8 bytes, a lightweight alternative for high-volume storage:

| Type | Size | Use Case |
|------|------|----------|
| `time.Time` | 24 bytes | General purpose |
| `tinytime.TinyTime` | ~8 bytes | High-volume events, logs, analytics |

## Key Concepts

1. **Unix vs Milliseconds**: This timestamp is in seconds (standard Unix format). Milliseconds would be ~1.5 trillion.
2. **Time zones**: Go defaults to UTC when parsing Unix timestamps.
3. **Precision**: TinyTime stores dates from year 0 to 9999, good for most applications.

## Common Patterns

```go
// Current time as timestamp
time.Now().Unix()        // seconds
time.Now().UnixMilli()   // milliseconds

// Convert back
time.Unix(1585750374, 0)
```

## When to Use TinyTime

- Storing many timestamps (database columns, cache keys)
- Memory-constrained environments
- When you don't need the full `time.Time` API
