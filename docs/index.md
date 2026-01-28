# toolruntime

`toolruntime` defines the runtime and trust boundary for executing code in
`toolcode`. It routes execution to backends with configurable security profiles.

## Key APIs

- `Runtime` interface
- `DefaultRuntime` implementation
- `Backend` interface
- `SecurityProfile` (dev/standard/hardened)
- `ExecuteRequest` / `ExecuteResult`

## Quickstart (dev)

```go
backend := unsafe.New(unsafe.Config{Mode: unsafe.ModeSubprocess})

rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
  Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
    toolruntime.ProfileDev: backend,
  },
  DefaultProfile: toolruntime.ProfileDev,
})
```

## Next

- Runtime pipeline and profiles: `architecture.md`
- Backend configuration: `usage.md`
- Examples: `examples.md`
