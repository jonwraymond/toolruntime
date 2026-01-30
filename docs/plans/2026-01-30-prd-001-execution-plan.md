# PRD-001 Execution Plan — toolruntime (TDD)

**Status:** Done
**Date:** 2026-01-30
**PRD:** `2026-01-30-prd-001-interface-contracts.md`


## TDD Workflow (required)
1. Red — write failing contract tests
2. Red verification — run tests
3. Green — minimal code/doc changes
4. Green verification — run tests
5. Commit — one commit per task


## Tasks
### Task 0 — Inventory + contract outline
- Confirm interface list and method signatures.
- Draft explicit contract bullets for each interface.
- Update docs/plans/README.md with this PRD + plan.
### Task 1 — Contract tests (Red/Green)
- Add `*_contract_test.go` with tests for each interface listed below.
- Use stub implementations where needed.
### Task 2 — GoDoc contracts
- Add/expand GoDoc on each interface with explicit contract clauses (thread-safety, errors, context, ownership).
- Update README/design-notes if user-facing.
### Task 3 — Verification
- Run `go test ./...`
- Run linters if configured (golangci-lint / gosec).


## Test Skeletons (contract_test.go)
### Backend
```go
func TestBackend_Contract(t *testing.T) {
    // Methods:
    // - Kind() BackendKind
    // - Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Runtime
```go
func TestRuntime_Contract(t *testing.T) {
    // Methods:
    // - Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### toolCallRecorder
```go
func TesttoolCallRecorder_Contract(t *testing.T) {
    // Methods:
    // - GetToolCalls() []ToolCallRecord
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### ToolGateway
```go
func TestToolGateway_Contract(t *testing.T) {
    // Methods:
    // - SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error)
    // - ListNamespaces(ctx context.Context) ([]string, error)
    // - DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error)
    // - ListToolExamples(ctx context.Context, id string, maxExamples int) ([]tooldocs.ToolExample, error)
    // - RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error)
    // - RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Connection
```go
func TestConnection_Contract(t *testing.T) {
    // Methods:
    // - Send(ctx context.Context, msg Message) error
    // - Receive(ctx context.Context) (Message, error)
    // - Close() error
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Codec
```go
func TestCodec_Contract(t *testing.T) {
    // Methods:
    // - Encode(msg Message) ([]byte, error)
    // - Decode(data []byte) (Message, error)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
### Logger
```go
func TestLogger_Contract(t *testing.T) {
    // Methods:
    // - Info(msg string, args ...any)
    // - Warn(msg string, args ...any)
    // - Error(msg string, args ...any)
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
