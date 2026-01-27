package toolruntime

import "context"

// Backend is the interface for code execution backends.
// Each backend provides a different level of isolation and security.
type Backend interface {
	// Kind returns the backend kind identifier.
	Kind() BackendKind

	// Execute runs code with the given request parameters.
	// It validates the request, executes the code, and returns the result.
	Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
}
