package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Category represents a stable error class with a documented exit code.
type Category string

const (
	CategoryAPI        Category = "api"
	CategoryAuth       Category = "auth"
	CategoryValidation Category = "validation"
	CategoryDiscovery  Category = "discovery"
	CategoryInternal   Category = "internal"
)

// Error is the structured error model for soke-cli.
type Error struct {
	Category   Category
	Message    string
	Operation  string
	ServerKey  string
	Retryable  bool
	Reason     string
	Hint       string
	Actions    []string
	Snapshot   string
	RPCCode    int               `json:"rpc_code,omitempty"`
	RPCData    json.RawMessage   `json:"rpc_data,omitempty"`
	ServerDiag ServerDiagnostics `json:"-"`
	Cause      error             `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying cause, enabling errors.Is and errors.As chains.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Option mutates a structured error before it is returned.
type Option func(*Error)

// ExitCode returns the documented process exit code for the error category.
func (e *Error) ExitCode() int {
	switch e.Category {
	case CategoryAPI:
		return 1
	case CategoryAuth:
		return 2
	case CategoryValidation:
		return 3
	case CategoryDiscovery:
		return 6
	default:
		return 5
	}
}

// WithOperation records the operation that failed.
func WithOperation(operation string) Option {
	return func(err *Error) {
		err.Operation = operation
	}
}

// WithServerKey records the server identifier associated with the failure.
func WithServerKey(serverKey string) Option {
	return func(err *Error) {
		err.ServerKey = serverKey
	}
}

// WithRetryable marks whether the error can be retried safely.
func WithRetryable(retryable bool) Option {
	return func(err *Error) {
		err.Retryable = retryable
	}
}

// WithReason records a stable machine-readable failure reason.
func WithReason(reason string) Option {
	return func(err *Error) {
		err.Reason = reason
	}
}

// WithHint records a short recovery hint for humans and agents.
func WithHint(hint string) Option {
	return func(err *Error) {
		err.Hint = hint
	}
}

// WithActions records suggested next actions for recovery.
func WithActions(actions ...string) Option {
	return func(err *Error) {
		out := make([]string, 0, len(actions))
		for _, action := range actions {
			if action == "" {
				continue
			}
			out = append(out, action)
		}
		if len(out) > 0 {
			err.Actions = out
		}
	}
}

// WithSnapshot records the recovery snapshot path associated with the failure.
func WithSnapshot(path string) Option {
	return func(err *Error) {
		err.Snapshot = path
	}
}

// WithRPCCode records the original JSON-RPC error code.
func WithRPCCode(code int) Option {
	return func(err *Error) {
		err.RPCCode = code
	}
}

// WithRPCData records the original JSON-RPC error data payload.
func WithRPCData(data json.RawMessage) Option {
	return func(err *Error) {
		err.RPCData = data
	}
}

// WithCause wraps the original error so it can be retrieved via errors.Unwrap.
func WithCause(err error) Option {
	return func(e *Error) {
		e.Cause = err
	}
}

func newError(category Category, message string, opts ...Option) error {
	err := &Error{
		Category: category,
		Message:  message,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(err)
		}
	}
	return err
}

// NewAPI returns an API-category error.
func NewAPI(message string, opts ...Option) error {
	return newError(CategoryAPI, message, opts...)
}

// NewAuth returns an auth-category error.
func NewAuth(message string, opts ...Option) error {
	return newError(CategoryAuth, message, opts...)
}

// NewValidation returns a validation-category error.
func NewValidation(message string, opts ...Option) error {
	return newError(CategoryValidation, message, opts...)
}

// NewDiscovery returns a discovery-category error.
func NewDiscovery(message string, opts ...Option) error {
	return newError(CategoryDiscovery, message, opts...)
}

// NewInternal returns an internal-category error.
func NewInternal(message string, opts ...Option) error {
	return newError(CategoryInternal, message, opts...)
}

// ExitCoder is implemented by errors that provide their own exit code.
type ExitCoder interface {
	ExitCode() int
}

// ExitCode maps any error to a stable exit code.
func ExitCode(err error) int {
	var typed *Error
	if errors.As(err, &typed) {
		return typed.ExitCode()
	}
	var ec ExitCoder
	if errors.As(err, &ec) {
		return ec.ExitCode()
	}
	return 5
}

// PrintJSON writes a machine-readable JSON error object.
func PrintJSON(w io.Writer, err error) error {
	errorPayload := map[string]any{
		"code":     ExitCode(err),
		"category": category(err),
		"message":  err.Error(),
	}

	var typed *Error
	if errors.As(err, &typed) {
		if typed.Reason != "" {
			errorPayload["reason"] = typed.Reason
		}
		if typed.Operation != "" {
			errorPayload["operation"] = typed.Operation
		}
		if typed.ServerKey != "" {
			errorPayload["server_key"] = typed.ServerKey
		}
		if typed.Retryable {
			errorPayload["retryable"] = true
		}
		if typed.Hint != "" {
			errorPayload["hint"] = typed.Hint
		}
		if len(typed.Actions) > 0 {
			errorPayload["actions"] = typed.Actions
		}
		if typed.Snapshot != "" {
			errorPayload["snapshot_path"] = typed.Snapshot
		}
		if typed.RPCCode != 0 {
			errorPayload["rpc_code"] = typed.RPCCode
		}
		if len(typed.RPCData) > 0 {
			var parsed any
			if json.Unmarshal(typed.RPCData, &parsed) == nil {
				errorPayload["rpc_data"] = parsed
			}
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(errorPayload)
}

// PrintHuman writes a human-readable error message to stderr.
func PrintHuman(w io.Writer, err error) error {
	if err == nil {
		return nil
	}

	lines := []string{fmt.Sprintf("Error: %s", err.Error())}

	var typed *Error
	if errors.As(err, &typed) {
		// Always show hint and actions when present
		if typed.Hint != "" {
			lines = append(lines, fmt.Sprintf("Hint: %s", typed.Hint))
		}
		if len(typed.Actions) > 0 {
			for _, action := range typed.Actions {
				if strings.TrimSpace(action) == "" {
					continue
				}
				lines = append(lines, fmt.Sprintf("Action: %s", action))
			}
		}
		if typed.Retryable {
			lines = append(lines, "Retryable: true")
		}
	}

	_, writeErr := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return writeErr
}

func category(err error) string {
	var typed *Error
	if errors.As(err, &typed) {
		return string(typed.Category)
	}
	return string(CategoryInternal)
}
