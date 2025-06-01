package errors

import (
	"errors"
	"fmt"

	errorscontracts "github.com/jeanmolossi/ai-agent-cli/app/contracts/errors"
)

type errorString struct {
	text   string
	module string
	args   []any
}

func New(text string, module ...string) errorscontracts.Error {
	err := &errorString{
		text: text,
	}

	if len(module) > 0 {
		err.module = module[0]
	}

	return err
}

// Args implements errors.Error.
func (e *errorString) Args(args ...any) errorscontracts.Error {
	e.args = args
	return e
}

// Error implements errors.Error.
func (e *errorString) Error() string {
	formattedText := e.text

	if len(e.args) > 0 {
		formattedText = fmt.Sprintf(e.text, e.args...)
	}

	if e.module != "" {
		formattedText = fmt.Sprintf("[%s] %s", e.module, formattedText)
	}

	return formattedText
}

// SetModule implements errors.Error.
func (e *errorString) SetModule(module string) errorscontracts.Error {
	e.module = module
	return e
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
