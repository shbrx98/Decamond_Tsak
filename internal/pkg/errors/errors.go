package errors

import "fmt"

type ValidationError struct{ msg string }
func (e *ValidationError) Error() string   { return e.msg }
func NewValidationError(msg string) error  { return &ValidationError{msg} }

type NotFoundError struct{ msg string }
func (e *NotFoundError) Error() string   { return e.msg }
func NewNotFoundError(msg string) error  { return &NotFoundError{msg} }

type UnauthorizedError struct{ msg string }
func (e *UnauthorizedError) Error() string   { return e.msg }
func NewUnauthorizedError(msg string) error  { return &UnauthorizedError{msg} }

type RateLimitError struct{ msg string }
func (e *RateLimitError) Error() string   { return e.msg }
func NewRateLimitError(msg string) error  { return &RateLimitError{msg} }

type InternalError struct{ msg string }
func (e *InternalError) Error() string   { return e.msg }
func NewInternalError(msg string) error  { return &InternalError{msg} }

// Helpers
func Wrap(err error, msg string) error {
    if err == nil { return nil }
    return fmt.Errorf("%s: %w", msg, err)
}