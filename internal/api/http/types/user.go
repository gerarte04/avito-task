package types

import (
	"avito-task/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

// Requests --------------------------------------------------

type SetIsActiveRequest struct {
	UserID string `json:"user_id"`
	IsActive bool `json:"is_active"`
}

func CreateSetIsActiveRequest(r *http.Request) (*SetIsActiveRequest, error) {
	const op = "CreateSetIsActiveRequest"

	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(req.UserID) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}

	return &req, nil
}

type GetReviewRequest struct {
	UserID string
}

func CreateGetReviewRequest(r *http.Request) (*GetReviewRequest, error) {
	const op = "CreateGetReviewRequest"

	req := GetReviewRequest{UserID: r.URL.Query().Get("user_id")}

	if len(req.UserID) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}

	return &req, nil
}

// Responses -------------------------------------------------

type SetIsActiveResponse struct {
	User *domain.User `json:"user"`
}

func CreateSetIsActiveResponse(user *domain.User) *SetIsActiveResponse {
	return &SetIsActiveResponse{User: user}
}

type GetReviewResponse struct {
	UserID string `json:"user_id"`
	PullRequests []*domain.PullRequestShort `json:"pull_requests"`
}

func CreateGetReviewResponse(id string, prs []*domain.PullRequestShort) *GetReviewResponse {
	return &GetReviewResponse{
		UserID: id,
		PullRequests: prs,
	}
}
