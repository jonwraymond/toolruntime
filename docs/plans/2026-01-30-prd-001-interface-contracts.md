# PRD-001 Interface Contracts — toolruntime

**Status:** Done
**Date:** 2026-01-30


## Overview
Define explicit interface contracts (GoDoc + documented semantics) for all interfaces in this repo. Contracts must state concurrency guarantees, error semantics, ownership of inputs/outputs, and context handling.


## Goals
- Every interface has explicit GoDoc describing behavioral contract.
- Contract behavior is codified in tests (contract tests).
- Docs/README updated where behavior is user-facing.


## Non-Goals
- No API shape changes unless required to satisfy the contract tests.
- No new features beyond contract clarity and tests.


## Interface Inventory
| Interface | File | Methods |
| --- | --- | --- |
| `Backend` | `toolruntime/backend.go:7` | Kind() BackendKind<br/>Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error) |
| `Runtime` | `toolruntime/runtime.go:12` | Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error) |
| `Logger` | `toolruntime/runtime.go:38` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `toolCallRecorder` | `toolruntime/runtime.go:56` | GetToolCalls() []ToolCallRecord |
| `ToolGateway` | `toolruntime/types.go:260` | SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error)<br/>ListNamespaces(ctx context.Context) ([]string, error)<br/>DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error)<br/>ListToolExamples(ctx context.Context, id string, maxExamples int) ([]tooldocs.ToolExample, error)<br/>RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error)<br/>RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) |
| `Connection` | `toolruntime/gateway/proxy/protocol.go:32` | Send(ctx context.Context, msg Message) error<br/>Receive(ctx context.Context) (Message, error)<br/>Close() error |
| `Codec` | `toolruntime/gateway/proxy/protocol.go:44` | Encode(msg Message) ([]byte, error)<br/>Decode(data []byte) (Message, error) |
| `Logger` | `toolruntime/backend/temporal/temporal.go:31` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/docker/docker.go:64` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/gvisor/gvisor.go:27` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/containerd/containerd.go:27` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/wasm/wasm.go:30` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/unsafe/unsafe.go:46` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/kata/kata.go:27` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/firecracker/firecracker.go:28` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/kubernetes/kubernetes.go:27` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |
| `Logger` | `toolruntime/backend/remote/remote.go:27` | Info(msg string, args ...any)<br/>Warn(msg string, args ...any)<br/>Error(msg string, args ...any) |

## Contract Template (apply per interface)
- **Thread-safety:** explicitly state if safe for concurrent use.
- **Context:** cancellation/deadline handling (if context is a parameter).
- **Errors:** classification, retryability, and wrapping expectations.
- **Ownership:** who owns/allocates inputs/outputs; mutation expectations.
- **Determinism/order:** ordering guarantees for returned slices/maps/streams.
- **Nil/zero handling:** behavior for nil inputs or empty values.


## Acceptance Criteria
- All interfaces have GoDoc with explicit behavioral contract.
- Contract tests exist and pass.
- No interface contract contradictions across repos.
