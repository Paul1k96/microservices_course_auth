package errs

// ErrorCode represents domain-level error codes.
type ErrorCode int

// Domain-level error codes.
const (
	Unknown ErrorCode = 1000 + iota
	UserNotFound
)

var (
	ErrUnknown      = NewError(Unknown, "unknown error")       // ErrUnknown represents an unknown error.
	ErrUserNotFound = NewError(UserNotFound, "user not found") // ErrUserNotFound represents a user not found error.
)

// Error represents an error.
type Error struct {
	Code    ErrorCode
	Message string
}

// NewError creates a new error.
func NewError(code ErrorCode, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}
