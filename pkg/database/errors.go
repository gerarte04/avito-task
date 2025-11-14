package database

import "errors"

var (
	ErrNotNullViolation = errors.New("not null violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrUniqueViolation = errors.New("unique violation")
	ErrUndocumented = errors.New("undocumented database error")
)
