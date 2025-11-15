package response

import (
	"avito-task/internal/repository"
	"errors"
	"log"
	"net/http"

	pkgErrors "avito-task/pkg/errors"
)

type ErrCodes struct {
	HTTPCode int
	StrCode  string
}

type ErrorDetails struct {
	StrCode string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Details ErrorDetails `json:"error"`
}

var (
	ErrInternal = errors.New("internal server error")
	ErrNotFound = errors.New("resource not found")

	errCodes = map[error]ErrCodes{
		ErrInternal: {http.StatusInternalServerError, "INTERNAL_ERROR"},
		ErrNotFound: {http.StatusNotFound, "NOT_FOUND"},

		repository.ErrTeamNameExists: {http.StatusBadRequest, "TEAM_EXISTS"},
		repository.ErrTeamNotExists:  {http.StatusNotFound, "NOT_FOUND"},
		repository.ErrUserNotExists:  {http.StatusNotFound, "NOT_FOUND"},

		repository.ErrPRIDExists:  {http.StatusConflict, "PR_EXISTS"},
		repository.ErrPRNotExists: {http.StatusNotFound, "NOT_FOUND"},
		repository.ErrPRMerged:    {http.StatusConflict, "PR_MERGED"},
		repository.ErrNotAssigned: {http.StatusConflict, "NOT_ASSIGNED"},
		repository.ErrNoCandidate: {http.StatusConflict, "NO_CANDIDATE"},
	}
)

func ProcessCreatingRequestError(w http.ResponseWriter, err error) {
	log.Print("[ERROR] ", err.Error())

	err = pkgErrors.UnwrapAll(err)

	WriteResponse(w, http.StatusBadRequest, ErrorResponse{
		Details: ErrorDetails{
			StrCode: "BAD_REQUEST",
			Message: err.Error(),
		},
	})
}

func ProcessError(w http.ResponseWriter, err error) {
	log.Print("[ERROR] ", err.Error())

	err = pkgErrors.UnwrapAll(err)
	codes := errCodes[ErrInternal]

	if docCode, ok := errCodes[err]; ok {
		codes = docCode
	} else {
		err = ErrInternal
	}

	if codes.HTTPCode == http.StatusNotFound {
		err = ErrNotFound
		codes = errCodes[ErrNotFound]
	}

	WriteResponse(w, codes.HTTPCode, ErrorResponse{
		Details: ErrorDetails{
			StrCode: codes.StrCode,
			Message: err.Error(),
		},
	})
}
