# Goroutine Leak Example

This example demonstrates a common goroutine leak pattern in Go.

## The Problem

The `gen()` function creates a goroutine that sends values to a channel in an infinite loop:

```go
func gen() <-chan int {
    ch := make(chan int)

    go func() {
        for i := 0; ; i++ {
            ch <- i
        }
    }()

    return ch
}
```

When `main()` breaks after receiving 5 values, the goroutine continues running and blocks forever on `ch <- i` waiting for a receiver that will never come.

## Why It's a Leak

- The goroutine is still alive after `main()` exits its loop
- No code can ever receive from the channel anymore
- The goroutine will block forever on the send operation
- This consumes memory and prevents garbage collection

## Detection

Run the program to see the leak:

```bash
go run main.go
```

You'll see:

- Number of goroutines > 1 (should be 1 after main finishes)
- pprof stack traces showing the blocked goroutine

## Solutions

### Option 1: Context Cancellation

```go
func gen(ctx context.Context) <-chan int {
    ch := make(chan int)

    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                return
            case ch <- i:
            }
        }
    }()

    return ch
}
```

### Option 2: Done Channel

```go
func gen(done <-chan struct{}) <-chan int {
    ch := make(chan int)

    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case <-done:
                return
            case ch <- i:
            }
        }
    }()

    return ch
}
```

## Key Takeaways

1. **Always provide cancellation signals** for goroutines you start
2. **Use `defer close(ch)`** to clean up channels
3. **Monitor goroutines** with `runtime.NumGoroutine()` and `pprof`
4. **Never create unbounded goroutines** without cleanup logic
