# Migration Guide: toolruntime to toolexec/runtime

This guide helps you migrate from `github.com/jonwraymond/toolruntime` to `github.com/jonwraymond/toolexec/runtime`.

## Import Path Changes

| Old Import | New Import |
|------------|------------|
| `github.com/jonwraymond/toolruntime` | `github.com/jonwraymond/toolexec/runtime` |
| `github.com/jonwraymond/toolruntime/backend` | `github.com/jonwraymond/toolexec/runtime/backend` |
| `github.com/jonwraymond/toolruntime/backend/unsafe` | `github.com/jonwraymond/toolexec/runtime/backend/unsafe` |
| `github.com/jonwraymond/toolruntime/backend/wasm` | `github.com/jonwraymond/toolexec/runtime/backend/wasm` |
| `github.com/jonwraymond/toolruntime/toolcodeengine` | `github.com/jonwraymond/toolexec/runtime/toolcodeengine` |

## Quick Migration Steps

1. **Update go.mod**

   Remove the old dependency:
   ```bash
   go mod edit -droprequire github.com/jonwraymond/toolruntime
   ```

   Add the new dependency:
   ```bash
   go get github.com/jonwraymond/toolexec@latest
   ```

2. **Update imports**

   Use `sed` or your editor to replace imports:
   ```bash
   find . -name "*.go" -exec sed -i '' 's|github.com/jonwraymond/toolruntime|github.com/jonwraymond/toolexec/runtime|g' {} +
   ```

3. **Run go mod tidy**
   ```bash
   go mod tidy
   ```

4. **Verify the build**
   ```bash
   go build ./...
   go test ./...
   ```

## API Compatibility

The API in `toolexec/runtime` is fully compatible with the former `toolruntime` package. All types, interfaces, and functions have been preserved:

- `Runtime` interface
- `Backend` interface
- `SecurityProfile` constants (`ProfileDev`, `ProfileStandard`, `ProfileHardened`)
- `ToolGateway` type
- `RuntimeConfig` struct
- `NewDefaultRuntime()` constructor

## Example: Before and After

### Before (toolruntime)

```go
import (
    "github.com/jonwraymond/toolruntime"
    "github.com/jonwraymond/toolruntime/backend/unsafe"
    "github.com/jonwraymond/toolruntime/toolcodeengine"
)

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

engine, err := toolcodeengine.New(toolcodeengine.Config{
    Runtime: rt,
    Profile: toolruntime.ProfileDev,
})
```

### After (toolexec/runtime)

```go
import (
    "github.com/jonwraymond/toolexec/runtime"
    "github.com/jonwraymond/toolexec/runtime/backend/unsafe"
    "github.com/jonwraymond/toolexec/runtime/toolcodeengine"
)

backend := unsafe.New(unsafe.Config{
    Mode:         unsafe.ModeSubprocess,
    RequireOptIn: false,
})

rt := runtime.NewDefaultRuntime(runtime.RuntimeConfig{
    Backends: map[runtime.SecurityProfile]runtime.Backend{
        runtime.ProfileDev: backend,
    },
    DefaultProfile: runtime.ProfileDev,
})

engine, err := toolcodeengine.New(toolcodeengine.Config{
    Runtime: rt,
    Profile: runtime.ProfileDev,
})
```

## Getting Help

If you encounter issues during migration, please open an issue in the [toolexec repository](https://github.com/jonwraymond/toolexec/issues).
