# Architecture

`toolruntime` is the execution boundary underneath `toolcode`.
It chooses a backend based on the requested security profile.

```mermaid
flowchart LR
  A[toolcode Engine] --> B[Runtime]
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

## Profiles

- `dev`: convenience, unsafe host allowed
- `standard`: safer defaults, may deny unsafe backend
- `hardened`: expected to require strong isolation
