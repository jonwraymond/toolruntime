# Usage

## Create a runtime

```go
rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
  Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
    toolruntime.ProfileDev: unsafe.New(unsafe.Config{Mode: unsafe.ModeSubprocess}),
  },
  DefaultProfile: toolruntime.ProfileDev,
})
```

## Execute code

```go
res, err := rt.Execute(ctx, toolruntime.ExecuteRequest{
  Language: "go",
  Code:     "__out = 2 + 2",
  Profile:  toolruntime.ProfileDev,
  Gateway:  gw,
})
```

## WASM backend (interface)

`toolruntime` defines the WASM backend interface in `backend/wasm`. You can
wire any compliant runtime (wazero/wasmtime/wasmer) by supplying a
`Runner` implementation.

```go
wasmBackend := wasm.New(wasm.Config{
  Runtime:    "wazero",
  EnableWASI: true,
  Client:     myWasmRunner, // implements wasm.Runner
})

rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
  Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
    toolruntime.ProfileStandard: wasmBackend,
  },
  DefaultProfile: toolruntime.ProfileStandard,
})
```

## Deny unsafe backend

```go
rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
  Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
    toolruntime.ProfileDev: unsafe.New(unsafe.Config{Mode: unsafe.ModeSubprocess}),
  },
  DenyUnsafeProfiles: []toolruntime.SecurityProfile{toolruntime.ProfileStandard},
  DefaultProfile:     toolruntime.ProfileStandard,
})
```
