# Examples

## Toolcode integration

```go
engine := toolcodeengine.New(toolcodeengine.Config{
  Runtime: rt,
  Profile: toolruntime.ProfileDev,
})

exec, _ := toolcode.NewDefaultExecutor(toolcode.Config{
  Index:  idx,
  Docs:   docs,
  Run:    runner,
  Engine: engine,
})
```

## Inspect tool calls

```go
res, _ := rt.Execute(ctx, toolruntime.ExecuteRequest{
  Language: "go",
  Code:     "__out = 1",
  Profile:  toolruntime.ProfileDev,
  Gateway:  gw,
})

for _, call := range res.ToolCalls {
  fmt.Println(call.ToolID, call.Duration)
}
```
