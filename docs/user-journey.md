# User Journey

This journey shows how `toolruntime` powers secure code execution in the broader stack.

## End-to-end flow (stack view)

![Diagram](assets/diagrams/user-journey.svg)

## Step-by-step

1. **Configure runtime** with backends keyed by security profile.
2. **Wrap tool access** using a `ToolGateway` (direct or proxy).
3. **Execute request** is routed to the backend based on profile.
4. **Backend runs code** with controlled access to tools via the gateway.
5. **Result and tool calls** are returned, with `LimitsEnforced` indicating actual enforcement.

## Example: wire runtime to toolcode

```go
rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
  Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
    toolruntime.ProfileStandard: mySandboxBackend,
  },
  DefaultProfile: toolruntime.ProfileStandard,
})

engine, err := toolcodeengine.New(toolcodeengine.Config{
  Runtime: rt,
  Profile: toolruntime.ProfileStandard,
})
if err != nil {
  return err
}

exec, _ := toolcode.NewDefaultExecutor(toolcode.Config{
  Index: idx,
  Docs:  docs,
  Run:   runner,
  Engine: engine,
})
```

## Expected outcomes

- Clear security posture by profile.
- Consistent tool access via a gateway boundary.
- Transparent limit enforcement reporting.

## Common failure modes

- `ErrRuntimeUnavailable` if no backend is registered for the profile.
- `ErrBackendDenied` when policy blocks unsafe backends.
- `ErrTimeout` / `ErrResourceLimit` on enforced limits.
