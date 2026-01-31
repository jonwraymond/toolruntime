# toolruntime

> **DEPRECATED**: This repository has been deprecated. The runtime functionality has been consolidated into the `toolexec/runtime` package.

## Migration

Please migrate to `github.com/jonwraymond/toolexec/runtime`. See [MIGRATION.md](MIGRATION.md) for detailed migration instructions.

## Why the change?

The `toolruntime` package has been merged into `toolexec` to simplify the dependency graph and provide a more cohesive execution layer. The `toolexec/runtime` package contains all the functionality previously provided by `toolruntime`, including:

- Backend-agnostic `Runtime` interface
- Security profiles (`dev`, `standard`, `hardened`)
- `ToolGateway` surface for safe tool access from sandboxes
- WASM backend interfaces

## Timeline

- This repository will remain available for reference but will not receive updates.
- New features and bug fixes will only be applied to `toolexec/runtime`.

## Documentation

For current documentation, see the [toolexec documentation](https://jonwraymond.github.io/ai-tools-stack/).
