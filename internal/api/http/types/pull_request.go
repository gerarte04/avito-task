package types

import (
	"avito-task/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

// Requests --------------------------------------------------

type CreatePRRequest struct {
	PR *domain.PullRequest
}

func MakeCreatePRRequest(r *http.Request) (*CreatePRRequest, error) {
	const op = "MakeCreatePRRequest"

	var pr domain.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(pr.ID) == 0 || len(pr.Name) == 0 || len(pr.AuthorID) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}

	return &CreatePRRequest{PR: &pr}, nil
}

type MergePRRequest struct {
	PRID string `json:"pull_request_id"`
}

func CreateMergePRRequest(r *http.Request) (*MergePRRequest, error) {
	const op = "CreateMergePRRequest"

	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(req.PRID) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}

	return &req, nil
}

type ReassignRequest struct {
	PRID string `json:"pull_request_id"`
	OldRewID string `json:"old_reviewer_id"`
}

func CreateReassignRequest(r *http.Request) (*ReassignRequest, error) {
	const op = "CreateReassignRequest"

	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(req.PRID) == 0 || len(req.OldRewID) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}

	return &req, nil
}

// Responses -------------------------------------------------

type CreatePRResponse struct {
	PR *domain.PullRequest `json:"pr"`
}

func MakeCreatePRResponse(pr *domain.PullRequest) *CreatePRResponse {
	return &CreatePRResponse{PR: pr}
}

type MergePRResponse struct {
	PR *domain.PullRequest `json:"pr"`
}

func CreateMergePRResponse(pr *domain.PullRequest) *MergePRResponse {
	return &MergePRResponse{PR: pr}
}

type ReassignResponse struct {
	ReplacedBy string `json:"replaced_by"`
	PR *domain.PullRequest `json:"pr"`
}

func CreateReassignResponse(newRewID string, pr *domain.PullRequest) *ReassignResponse {
	return &ReassignResponse{
		ReplacedBy: newRewID,
		PR: pr,
	}
}
