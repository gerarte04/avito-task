package repository

import "errors"

var (
	ErrTeamNotExists = errors.New("team not exists")
	ErrUserNotExists = errors.New("user not exists")
	ErrPRNotExists = errors.New("PR not exists")
)
