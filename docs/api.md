# API Reference

## Runtime

```go
type Runtime interface {
  Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
}
```

## Backend

```go
type Backend interface {
  Kind() BackendKind
  Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
}
```

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
  SearchTools(query string, limit int) ([]toolindex.Summary, error)
  DescribeTool(id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error)
  RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error)
  RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error)
}
```

## RuntimeConfig

```go
type RuntimeConfig struct {
  Backends           map[SecurityProfile]Backend
  DenyUnsafeProfiles []SecurityProfile
  DefaultProfile     SecurityProfile
  Logger             Logger
}
```
