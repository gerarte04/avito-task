package repository

import "errors"

var (
	ErrTeamNameExists = errors.New("team_name already exists")
	ErrTeamNotExists = errors.New("team not exists")

	ErrUserNotExists = errors.New("user not exists")

	ErrPRIDExists = errors.New("PR id already exists")
	ErrPRNotExists = errors.New("PR not exists")
	ErrPRMerged = errors.New("cannot reassign on merged PR")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
)
