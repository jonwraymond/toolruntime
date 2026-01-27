# toolruntime

`toolruntime` is the execution runtime and trust boundary underneath
`toolcode`. It provides:

- a backend-agnostic `Runtime` interface,
- security profiles (`dev`, `standard`, `hardened`), and
- a `ToolGateway` surface for safe tool access from sandboxes.

It is designed to plug into:

- `toolcode` via `toolruntime/toolcodeengine`, and
- `metatools-mcp` as the runtime behind `execute_code`.

Important current status:

- The `unsafe_host` backend works for dev, but it is not sandboxed.
- Most other backends are scaffolds with policy shaping, not full isolation.
- Snippet-to-tool wiring is still primarily enforced by `toolcode` today.

## Core concepts

- `Runtime`: routes an `ExecuteRequest` to the chosen backend.
- `Backend`: the isolation mechanism (unsafe host, docker, kubernetes, etc).
- `ToolGateway`: the only allowed tool surface for untrusted code.
- `toolcodeengine`: adapter that implements `toolcode.Engine`.

## Quickstart (dev only, unsafe host backend)

This shows the intended wiring: tool libs -> runtime -> toolcode engine ->
toolcode executor.

```go
package main

import (
  "context"
  "time"

  "github.com/jonwraymond/toolcode"
  "github.com/jonwraymond/tooldocs"
  "github.com/jonwraymond/toolindex"
  "github.com/jonwraymond/toolrun"
  "github.com/jonwraymond/toolruntime"
  "github.com/jonwraymond/toolruntime/backend/unsafe"
  "github.com/jonwraymond/toolruntime/toolcodeengine"
)

func main() {
  idx := toolindex.NewInMemoryIndex()
  docs := tooldocs.NewInMemoryStore(tooldocs.StoreOptions{Index: idx})
  runner := toolrun.NewRunner(toolrun.WithIndex(idx))

  backend := unsafe.New(unsafe.Config{
    Mode:         unsafe.ModeSubprocess,
    RequireOptIn: false,
  })

  rt := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
    Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
      toolruntime.ProfileDev: backend,
    },
    DefaultProfile: toolruntime.ProfileDev,
  })

  engine := toolcodeengine.New(toolcodeengine.Config{
    Runtime: rt,
    Profile: toolruntime.ProfileDev,
  })

  exec, err := toolcode.NewDefaultExecutor(toolcode.Config{
    Index:          idx,
    Docs:           docs,
    Run:            runner,
    Engine:         engine,
    DefaultTimeout: 10 * time.Second,
    MaxToolCalls:   64,
    MaxChainSteps:  8,
  })
  if err != nil {
    panic(err)
  }

  _, _ = exec.ExecuteCode(context.Background(), toolcode.ExecuteParams{
    Language: "go",
    Code:     `__out = "ok"`,
  })
}
```

## Security posture guidance

- `ProfileDev` + `unsafe_host` is for local development only.
- Treat all schemas/docs/annotations as untrusted input.
- For production, plan on container isolation and then stronger runtimes
  (gVisor/Kata/microVM).

## Version compatibility (current tags)

- `toolmodel`: `v0.1.0`
- `toolindex`: `v0.1.1`
- `tooldocs`: `v0.1.1`
- `toolrun`: `v0.1.0`
- `toolcode`: `v0.1.0`
- `toolruntime`: `v0.1.0`
- `metatools-mcp`: `v0.1.2`

## CI

CI runs:

- `go mod download`
- `go vet ./...`
- `go test ./...`
