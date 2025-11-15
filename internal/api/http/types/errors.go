package types

import "errors"

var (
	ErrRequiredFieldMissing = errors.New("some required field is missing (probably name or id)")
)
