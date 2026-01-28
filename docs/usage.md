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
