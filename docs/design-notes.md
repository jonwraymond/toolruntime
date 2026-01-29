# Design Notes

This page documents the tradeoffs and error semantics behind `toolruntime`.

## Design tradeoffs

- **Runtime as router.** `toolruntime` does not execute code itself; it routes requests to backends based on a security profile. This keeps policy separate from execution.
- **Explicit security profiles.** `dev`, `standard`, and `hardened` provide predictable isolation tiers. This trades flexibility for clarity and safer defaults.
- **Gateway boundary.** Executed code never talks to toolindex/tooldocs/toolrun directly. Instead it receives a `ToolGateway`, which preserves the trust boundary.
- **Context propagation.** Gateway operations check `ctx.Err()` and cancel early, so timeouts and cancellations flow through to tool discovery and execution calls.
- **Backend diversity.** Backends represent different isolation levels (host, containers, microVMs, WASM, remote). This lets deployments match security and performance needs.
- **Limits are declarative.** `Limits` are requested by the caller, and backends report which limits were actually enforced.

## Error semantics

`toolruntime` uses sentinel errors for classification:

- `ErrRuntimeUnavailable` – no backend exists for the requested profile.
- `ErrBackendDenied` – policy denies the backend (e.g., unsafe host for hardened profile).
- `ErrSandboxViolation` – sandbox policy was violated.
- `ErrTimeout` – execution exceeded timeout.
- `ErrResourceLimit` – resource limits exceeded.
- `ErrMissingGateway` / `ErrMissingCode` – invalid `ExecuteRequest`.
- `ErrInvalidLimits` – limits validation failed.

`RuntimeError` wraps backend failures with `BackendKind`, `Op`, and `Retryable`.

## Extension points

- **Custom backends:** implement `Backend` to integrate Docker, containerd, Kubernetes, gVisor, WASM, or a remote execution service.
- **Custom gateways:** use `gateway/direct` for in-process execution or `gateway/proxy` for RPC-mediated execution.
- **Toolcode integration:** `toolcodeengine` adapts `toolruntime` into a `toolcode.Engine`.

## Operational guidance

- Use `ProfileStandard` by default; reserve `ProfileDev` for trusted environments.
- Enforce `DenyUnsafeProfiles` to prevent accidental unsafe execution.
- Report `LimitsEnforced` accurately so callers can detect degraded enforcement.
