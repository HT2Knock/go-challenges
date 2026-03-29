# Goroutine Memory Measurement

## Two Approaches Demonstrated

### 1. runtime.MemStats

Direct memory statistics from the Go runtime.

```go
var m1, m2 runtime.MemStats
runtime.ReadMemStats(&m1)
// ... create goroutines ...
runtime.GC()
runtime.ReadMemStats(&m2)

memPerGoroutine := (m2.StackInuse - m1.StackInuse) / numGoroutines
```

**What it measures:**
- `StackInuse`: Memory currently used by goroutine stacks
- Initial goroutine stack: 2KB
- Grows dynamically as needed

**Typical results:**
- Per goroutine: ~2-8 KB (depends on stack usage)
- 10,000 goroutines: ~20-80 MB

### 2. pprof Heap Profile

Generates a file for analysis with `go tool pprof`.

```bash
go tool pprof goroutine.prof
```

**What you can see:**
- Goroutine count over time
- Stack traces showing where they were created
- Memory allocation by function

**Useful commands in pprof:**
```
top    - top goroutine consumers
list   - source code with goroutine counts
web    - open flame graph in browser
```

## When to Use Each

| Method | Best For |
|--------|----------|
| MemStats | Quick measurement, CI checks |
| pprof | Deep analysis, finding leak sources |

## Gotchas

1. **GC variability**: Go's garbage collector runs asynchronously. Call `runtime.GC()` before measuring for accurate results.
2. **Stack growth**: Goroutine stacks start small (2KB) and grow on demand. Sleeping goroutines use minimal memory.
3. **pprof overhead**: Writing profiles has overhead - don't use in production hot paths without sampling.

## Key Takeaway

Understanding goroutine memory helps diagnose:
- Memory leaks from abandoned goroutines
- Resource exhaustion from unbounded goroutine creation
- Performance issues from excessive context switching
