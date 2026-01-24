---
name: writing-go
description: Go coding conventions and idioms. Use when writing, reviewing, or generating Go code in *.go files.
compatibility:
  - opencode
  - claude
license: MIT
---

# Writing Go

## Function Signatures

### Context First
Context should be the first parameter and named `ctx`.

```go
// Good
func ProcessOrder(ctx context.Context, orderID string) error
func (s *Service) FetchUser(ctx context.Context, id int64) (*User, error)

// Bad
func ProcessOrder(orderID string, ctx context.Context) error
func (s *Service) FetchUser(id int64, c context.Context) (*User, error)
```

### Receiver Naming
Use short (1-2 letter) receiver names based on the type name. Never use `this` or `self`.

```go
// Good
func (s *Server) Start() error
func (s *Server) Stop() error
func (uc *UserController) GetUser(id int64) (*User, error)

// Bad
func (this *Server) Start() error
func (self *Server) Stop() error
func (server *Server) Restart() error
```

### No Naked Returns
Always specify return values explicitly. Do not use naked returns.

```go
// Good
func ParseConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, err
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, err
    }
    return cfg, nil
}

// Bad
func ParseConfig(path string) (cfg Config, err error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return
    }
    if err = json.Unmarshal(data, &cfg); err != nil {
        return
    }
    return
}
```

## Interfaces

### Interface Naming
Single-method interfaces use `-er` suffix. Multi-method interfaces use descriptive names.

```go
// Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Validator interface {
    Validate() error
}

type UserRepository interface {
    FindUser(id int64) (*User, error)
    SaveUser(u *User) error
}

// Bad
type IReader interface {
    Read(p []byte) (n int, err error)
}

type Validation interface {
    Validate() error
}

type UserFinderInterface interface {
    FindUser(id int64) (*User, error)
}
```

### No Unexported Returns
Do not return unexported types from exported functions.

```go
// Good
type User struct {
    ID   int64
    Name string
}

func NewUser(name string) *User

// Bad
type user struct {
    ID   int64
    Name string
}

func NewUser(name string) *user
```

## Error Handling

### Error Assignment
Use plain assignment for error handling, not inline declaration in if statements.

```go
// Good
err := doSomething()
if err != nil {
    return err
}

// Bad
if err := doSomething(); err != nil {
    return err
}
```

### Defer Close
Defer `Close()` immediately after error check, not before.

```go
// Good
f, err := os.Open(path)
if err != nil {
    return nil, err
}
defer f.Close()

// Bad
f, err := os.Open(path)
defer f.Close() // panic if f is nil
if err != nil {
    return nil, err
}
```

## Struct Tags

Use consistent struct tag formatting with proper casing and spacing.

```go
// Good
type User struct {
    ID        int64     `json:"id" db:"id"`
    FirstName string    `json:"first_name" db:"first_name"`
    Email     string    `json:"email,omitempty"`
    Internal  string    `json:"-"`
}

// Bad
type User struct {
    ID        int64  `json:"id"  db:"id"`
    FirstName string `json:"firstName" db:"first_name"`
    Email     string `JSON:"email"`
}
```

## Logging (slog)

### Use slog
Use `log/slog` for logging instead of `log` or `fmt.Print`.

```go
// Good
slog.Info("server started", "port", port)
slog.Error("request failed", "error", err)

// Bad
log.Printf("server started on port %d", port)
fmt.Println("request failed:", err)
```

### Context-Aware Logging
Use context-aware slog functions (`InfoContext`, `ErrorContext`, etc.) when context is available.

```go
// Good
func ProcessOrder(ctx context.Context, orderID string) error {
    slog.InfoContext(ctx, "processing order", slog.String("order_id", orderID))
    // ...
    slog.ErrorContext(ctx, "processing failed",
        slog.String("order_id", orderID),
        slog.Any("error", err))
    return err
}

// Bad
func ProcessOrder(ctx context.Context, orderID string) error {
    slog.Info("processing order", slog.String("order_id", orderID))  // ctx available but not used
    // ...
    slog.Error("processing failed",
        slog.String("order_id", orderID),
        slog.Any("error", err))
    return err
}
```

### Structured Logging
Use type-safe attribute constructors for better performance and type checking.

```go
// Good
slog.Info("user logged in",
    slog.String("user_id", userID),
    slog.String("ip", remoteAddr),
    slog.Int("attempt", attemptCount))

slog.Error("request failed",
    slog.String("path", r.URL.Path),
    slog.Duration("latency", elapsed),
    slog.Any("error", err))

// Bad
slog.Info(fmt.Sprintf("user %s logged in from %s", userID, remoteAddr))
slog.Info("user logged in", "user_id", userID, "ip", remoteAddr)  // untyped
slog.Error("request failed: " + err.Error())
```
