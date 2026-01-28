# toolruntime

Runtime/sandbox abstraction with multiple backends.

## What this repo provides

- Backends for unsafe, docker, kubernetes, gvisor, firecracker, wasm
- Profiles and execution limits
- Standardized runtime interface

## Example

```go
rt := toolruntime.NewDefaultRuntime(cfg)
```
