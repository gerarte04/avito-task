package types

import (
	"avito-task/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

// Requests --------------------------------------------------

type AddTeamRequest struct {
	Team *domain.Team
}

func CreateAddTeamRequest(r *http.Request) (*AddTeamRequest, error) {
	const op = "CreateAddTeamRequest"
	
	var team domain.Team
	
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	if len(team.Name) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}
	
	for _, u := range team.Members {
		if len(u.ID) == 0 || len(u.Name) == 0 {
			return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
		}
	}
	
	return &AddTeamRequest{Team: &team}, nil
}

type GetTeamRequest struct {
	Name string
}

func CreateGetTeamRequest(r *http.Request) (*GetTeamRequest, error) {
	const op = "CreateGetTeamRequest"
	
	var req GetTeamRequest
	req.Name = r.URL.Query().Get("team_name")
	
	if len(req.Name) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredFieldMissing)
	}
	
	return &req, nil
}

// Responses -------------------------------------------------

type AddTeamResponse struct {
	Team *domain.Team `json:"team"`
}

func CreateAddTeamResponse(team *domain.Team) *AddTeamResponse {
	for _, u := range team.Members {
		u.TeamName = ""
	}

	return &AddTeamResponse{Team: team}
}

func CreateGetTeamResponse(team *domain.Team) *domain.Team {
	for _, u := range team.Members {
		u.TeamName = ""
	}

	return team
}
