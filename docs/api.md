# API Reference

## Runtime

```go
type Runtime interface {
  Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
}
```

### Runtime contract

- Concurrency: implementations are safe for concurrent use.
- Context: honors cancellation/deadlines and returns `ctx.Err()` when canceled.
- Errors: request validation uses `ErrMissingGateway`/`ErrInvalidRequest`;
  backend selection uses `ErrRuntimeUnavailable`/`ErrBackendDenied`.

## Backend

```go
type Backend interface {
  Kind() BackendKind
  Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
}
```

### Backend contract

- Concurrency: implementations are safe for concurrent use.
- Context: honors cancellation/deadlines and returns `ctx.Err()` when canceled.
- Errors: validation should return `ErrInvalidRequest` where applicable.

## SecurityProfile

```go
type SecurityProfile string
const (
  ProfileDev      SecurityProfile = "dev"
  ProfileStandard SecurityProfile = "standard"
  ProfileHardened SecurityProfile = "hardened"
)
```

## ExecuteRequest / ExecuteResult

```go
type ExecuteRequest struct {
  Language string
  Code     string
  Args     map[string]any
  Profile  SecurityProfile
  Gateway  ToolGateway
  Timeout  time.Duration
}

type ExecuteResult struct {
  Value      any
  Stdout     string
  Stderr     string
  ToolCalls  []ToolCallRecord
  Duration   time.Duration
  Backend    BackendInfo
}
```

## ToolGateway

```go
type ToolGateway interface {
  SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error)
  ListNamespaces(ctx context.Context) ([]string, error)
  DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error)
  ListToolExamples(ctx context.Context, id string, maxExamples int) ([]tooldocs.ToolExample, error)
  RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error)
  RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error)
}
```

### ToolGateway contract

- Concurrency: implementations are safe for concurrent use.
- Context: honors cancellation/deadlines and returns `ctx.Err()` when canceled.
- Ownership: args are read-only; results are caller-owned snapshots.

## RuntimeConfig

```go
type RuntimeConfig struct {
  Backends           map[SecurityProfile]Backend
  DenyUnsafeProfiles []SecurityProfile
  DefaultProfile     SecurityProfile
  Logger             Logger
}
```

### Errors

- `ErrMissingGateway`
- `ErrInvalidRequest`
- `ErrRuntimeUnavailable`
- `ErrBackendDenied`
