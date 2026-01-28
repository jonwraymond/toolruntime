# User Journey

This journey shows how `toolruntime` powers secure code execution in the broader stack.

## End-to-end flow (stack view)

```mermaid
sequenceDiagram
  participant Agent
  participant MCP as metatools-mcp
  participant Code as toolcode
  participant Runtime as toolruntime
  participant Backend as Sandbox Backend
  participant Gateway as ToolGateway
  participant Runner as toolrun

  Agent->>MCP: execute_code
  MCP->>Code: ExecuteCode
  Code->>Runtime: Execute(profile, code)
  Runtime->>Backend: Execute
  Backend->>Gateway: RunTool/RunChain
  Gateway->>Runner: Run
  Runner-->>Backend: Result
  Backend-->>Runtime: ExecuteResult
  Runtime-->>Code: ExecuteResult
  Code-->>MCP: Output
```

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

engine := toolcodeengine.New(toolcodeengine.Config{
  Runtime: rt,
  Profile: toolruntime.ProfileStandard,
})

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
