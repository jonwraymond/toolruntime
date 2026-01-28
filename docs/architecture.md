# Architecture

`toolruntime` is the execution boundary underneath `toolcode`.
It chooses a backend based on the requested security profile.

## Runtime selection

```mermaid
flowchart LR
  A[ExecuteRequest] --> B[DefaultRuntime]
  B --> C[Profile selection]
  C --> D[Backend]

  subgraph Backends
    U[unsafe host]
    Dk[docker]
    K8s[kubernetes]
    Gv[gvisor]
    Fc[firecracker]
    W[wasm]
  end
```

## Execution sequence

```mermaid
sequenceDiagram
  participant Client
  participant RT as Runtime
  participant BE as Backend
  participant GW as ToolGateway

  Client->>RT: Execute(req)
  RT->>BE: Execute(req)
  BE->>GW: tool calls
  GW-->>BE: tool results
  BE-->>RT: ExecuteResult
  RT-->>Client: ExecuteResult
```

## Profiles

- `dev`: convenience, unsafe host allowed
- `standard`: safer defaults, may deny unsafe backend
- `hardened`: expected to require strong isolation
