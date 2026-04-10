# Deadlock Challenge

## What Happens

Two goroutines block forever, waiting on each other. The program hangs with no
output beyond the first two print statements.

```text
A: Trying to send ch1...
B: Trying to send ch2...
(hangs forever)
```

## Why This Happens

**Unbuffered channels** require both sender and receiver to be ready simultaneously. When:

- Goroutine A sends on ch1, it blocks waiting for someone to receive
- Goroutine B sends on ch2, it blocks waiting for someone to receive
- Now both are blocked - neither can receive from the other

## The Coffman Conditions (4 conditions for deadlock)

1. **Mutual Exclusion** - Channels can only be used by one goroutine at a time
2. **Hold and Wait** - Goroutines hold their send while waiting for the other channel
3. **No Preemption** - Can't force a goroutine to give up its channel operation
4. **Circular Wait** - A waits for B, B waits for A (the cycle)

## Fixes

### 1. Buffered Channels

```go
ch1 := make(chan struct{}, 1)  // capacity of 1
ch2 := make(chan struct{}, 1)
```

One goroutine can send immediately, the receiver will be ready by the time the other tries to send.

### 2. Select with Default (non-blocking)

```go
select {
case ch1 <- struct{}{}:
default:  // would block, do something else
}
```

### 3. Single Channel

If both goroutines need to synchronize, use one channel:

```go
ch := make(chan struct{})
```

## Key Takeaway

Always consider the synchronization order. Circular dependencies in channel communication lead to deadlock.
