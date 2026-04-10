# Agent Guidelines for go-challenges

Repository of Go learning challenges covering concurrency patterns, memory
management, and data structures.

## Project Structure

```text
go-challenges/
├── go.work              # Workspace file (includes all modules)
├── deadlock/            # Deadlock demonstration
├── goroutine-leak/      # Goroutine leak patterns
├── goroutine-memory/    # Memory measurement techniques
├── task-scheduling/      # Task graph scheduling (WIP)
└── time-conversion/     # Time library examples
```

## Build Commands

### Build All Modules

```bash
cd <module-directory> && go build -o /dev/null
```

### Run a Module

```bash
cd <module-directory> && go run .
```

### Build Specific Module (from repo root)

```bash
cd deadlock && go build
```

## Test Commands

### Run All Tests (per module)

```bash
cd <module-directory> && go test ./...
```

### Run a Single Test

```bash
cd <module-directory> && go test -v -run TestName
```

### Run Tests with Coverage

```bash
cd <module-directory> && go test -cover
```

### Example: Run specific test in goroutine-leak

```bash
cd goroutine-leak && go test -v -run Test_measureGoroutineMemory
```

## Code Style Guidelines

### Formatting

- Use `gofmt` or `go fmt` for formatting
- Indent with tabs, not spaces
- No trailing whitespace
- One blank line between import groups

### Imports

- Group stdlib and external imports
- Use `goimports` for automatic import management
- Order: stdlib → external → internal (if any)

```go
import (
    "fmt"
    "os"

    "github.com/example/pkg"
)
```

### Naming Conventions

- **Packages**: lowercase, no underscores (e.g., `deadlock`, not `dead_lock`)
- **Functions**: PascalCase (e.g., `GenerateTasks`)
- **Variables**: camelCase (e.g., `numGoroutines`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interfaces**: PascalCase, often with `-er` suffix (e.g., `Reader`, `Writer`)

### Error Handling

- Return errors as last return value: `(result, error)`
- Wrap errors with context: `fmt.Errorf("operation: %w", err)`
- Never ignore errors with `_`
- Use sentinel errors for known conditions: `var ErrNotFound = errors.New("not found")`

### Goroutines and Channels

- Always provide a way to signal goroutine completion (done channel, context)
- Close channels from the sender side only
- Use `sync.WaitGroup` for waiting on multiple goroutines
- Document concurrent behavior in comments

### Types

- Use concrete types unless interface is needed for polymorphism
- Prefer struct composition over inheritance
- Use `time.Duration` for time intervals, `time.Time` for timestamps
- Use `context.Context` for cancellation signals

### Comments

- Comment exported identifiers (will appear in godoc)
- Use `//` for single-line comments
- No commented-out code (delete it, use git history)
- Explain "why" not "what"

## Testing Guidelines

### Test Function Naming

```go
func Test_FunctionName_Behavior(t *testing.T) { }
```

### Table-Driven Tests

```go
tests := []struct {
    name    string
    input   string
    want    string
}{
    {"case1", "input1", "want1"},
    {"case2", "input2", "want2"},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        if got := fn(tt.input); got != tt.want {
            t.Errorf("fn() = %v, want %v", got, tt.want)
        }
    })
}
```

### Assertions

- Use `t.Errorf` or `t.Fatalf` for failures
- Include actual vs expected in error messages
- Use `require` package for fatal assertions (stops test immediately)

## Notes for This Repository

- Each module is independent with its own `go.mod`
- `go.work` references all modules for workspace support
- Use `go-tinytime` for lightweight timestamp storage (time-conversion module)
- Goroutine leak detection: use `runtime.NumGoroutine()` and `pprof`
- Deadlock demo uses unbuffered channels - see NOTES.md for fixes
