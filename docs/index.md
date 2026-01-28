# toolruntime

`toolruntime` defines the runtime and trust boundary for executing code in
`toolcode`. It routes execution to backends with configurable security profiles.

## Motivation

- **Isolation** for untrusted or semi-trusted code
- **Portability** across environments (host, containers, VMs)
- **Policy** via explicit security profiles

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

## Usability notes

- Profiles keep policy separate from implementation
- Backends are swappable without changing the executor API
- Tool calls are mediated via the gateway interface

## Next

- Runtime pipeline and profiles: `architecture.md`
- Backend configuration: `usage.md`
- Examples: `examples.md`
- Design Notes: `design-notes.md`
- User Journey: `user-journey.md`

!!! note "Deep dives"
    - Design Notes: `design-notes.md`
    - User Journey: `user-journey.md`
