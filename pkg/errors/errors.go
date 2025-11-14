package errors

import "errors"

func UnwrapAll(err error) error {
	var basicErr error

	for next := err; next != nil; next = errors.Unwrap(basicErr) {
		basicErr = next
	}

	return basicErr
}
