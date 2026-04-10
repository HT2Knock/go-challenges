# Goroutine Leak Challenge

## What Happens

A goroutine is spawned to generate integers infinitely. When the consumer breaks after reading 6 values, the producer goroutine keeps running forever.

```text
0
1
2
3
4
5
Number of Go routines: 2
```

Notice: Even though we "finished", there are still 2 goroutines (main + the generator).

## Why This Matters

1. **Memory leak**: Goroutines have initial 2KB stack, can grow to ~1GB
2. **Resource exhaustion**: 10,000 leaked goroutines = ~20MB just for stacks
3. **Production impact**: Handlers that spawn goroutines without cleanup will accumulate over time

## Detection Methods

### 1. runtime.NumGoroutine()

```go
fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
```

Quick sanity check during development.

### 2. pprof Goroutine Profile

```go
pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
```

Dumps all goroutines with stack traces - shows where leaks originate.

### 3. go tool pprof

```bash
# Generate profile
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Or from a saved profile
go tool pprof goroutine.prof
```

## Fixes

### 1. Done Channel

```go
func gen(done <-chan struct{}) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case ch <- i:
            case <-done:
                return
            }
        }
    }()
    return ch
}

// Usage
done := make(chan struct{})
defer close(done)
for i := range gen(done) {
    if i == 5 {
        break
    }
}
```

### 2. Context Cancellation

```go
func gen(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}

// Usage
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
for i := range gen(ctx) {
    // ...
}
```

## Key Takeaway

Always provide a way to signal goroutines to stop. The garbage collector does NOT automatically clean up running goroutines.
