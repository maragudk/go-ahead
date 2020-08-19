package errors2

import (
	"fmt"
)

// Wrap an error with a message, or return nil if err is nil.
func Wrap(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	var fmtArgs []interface{}
	fmtArgs = append(fmtArgs, args...)
	fmtArgs = append(fmtArgs, err)
	return fmt.Errorf(message+": %w", fmtArgs...)
}
