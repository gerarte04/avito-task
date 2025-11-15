package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	tu "avito-task/pkg/testutils"

	"github.com/stretchr/testify/require"
)

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func TestReviewerServiceIntegration(t *testing.T) {
	require := require.New(t)

	url := fmt.Sprintf("http://%s", os.Getenv("HTTP_ADDRESS"))
	resHC, _ := tu.MakeRequest(t, url, "GET", "/health", nil)
	require.Equal(http.StatusOK, resHC.StatusCode)

	var teamResponse struct {
		Team Team `json:"team"`
	}
	var userResponse struct {
		User User `json:"user"`
	}
	var prResponse struct {
		PR PullRequest `json:"pr"`
	}
	var reassignResponse struct {
		PR         PullRequest `json:"pr"`
		ReplacedBy string      `json:"replaced_by"`
	}
	var reviewListResponse struct {
		UserID       string             `json:"user_id"`
		PullRequests []PullRequestShort `json:"pull_requests"`
	}
	var errResponse ErrorResponse

	var createdTeam Team
	var createdPR PullRequest
	var originalReviewer string

	teamPayload := Team{
		TeamName: "backend-devs",
		Members: []TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
			{UserID: "u3", Username: "Charlie", IsActive: true},
			{UserID: "u4", Username: "David", IsActive: true},
		},
	}

	t.Run("A_TeamWorkflow", func(t *testing.T) {
		t.Run("1_CreateTeam_Success", func(t *testing.T) {
			res, body := tu.MakeRequest(t, url, "POST", "/team/add", teamPayload)
			require.Equal(http.StatusCreated, res.StatusCode)

			err := json.Unmarshal([]byte(body), &teamResponse)
			require.NoError(err)
			require.Equal(teamPayload.TeamName, teamResponse.Team.TeamName)
			require.Len(teamResponse.Team.Members, 4)
			createdTeam = teamResponse.Team
		})

		t.Run("2_CreateTeam_Duplicate", func(t *testing.T) {
			res, body := tu.MakeRequest(t, url, "POST", "/team/add", teamPayload)
			require.Equal(http.StatusBadRequest, res.StatusCode)

			err := json.Unmarshal([]byte(body), &errResponse)
			require.NoError(err)
			require.Equal("TEAM_EXISTS", errResponse.Error.Code)
		})

		t.Run("3_GetTeam_Success", func(t *testing.T) {
			path := fmt.Sprintf("/team/get?team_name=%s", createdTeam.TeamName)
			res, body := tu.MakeRequest(t, url, "GET", path, nil)
			require.Equal(http.StatusOK, res.StatusCode)

			var getTeamResponse Team
			err := json.Unmarshal([]byte(body), &getTeamResponse)
			require.NoError(err)
			require.Equal(createdTeam.TeamName, getTeamResponse.TeamName)
			require.Len(getTeamResponse.Members, 4)
		})

		t.Run("4_GetTeam_NotFound", func(t *testing.T) {
			res, _ := tu.MakeRequest(t, url, "GET", "/team/get?team_name=frontend", nil)
			require.Equal(http.StatusNotFound, res.StatusCode)
		})
	})

	t.Run("B_UserWorkflow", func(t *testing.T) {
		t.Run("1_SetInactive_Success", func(t *testing.T) {
			payload := map[string]interface{}{
				"user_id":   "u3",
				"is_active": false,
			}
			res, body := tu.MakeRequest(t, url, "POST", "/users/setIsActive", payload)
			require.Equal(http.StatusOK, res.StatusCode)

			err := json.Unmarshal([]byte(body), &userResponse)
			require.NoError(err)
			require.Equal("u3", userResponse.User.UserID)
			require.False(userResponse.User.IsActive)
			require.Equal(createdTeam.TeamName, userResponse.User.TeamName)
		})

		t.Run("2_SetInactive_NotFound", func(t *testing.T) {
			payload := map[string]interface{}{
				"user_id":   "u99",
				"is_active": false,
			}
			res, _ := tu.MakeRequest(t, url, "POST", "/users/setIsActive", payload)
			require.Equal(http.StatusNotFound, res.StatusCode)
		})
	})

	t.Run("C_PullRequestWorkflow", func(t *testing.T) {
		prPayload := map[string]string{
			"pull_request_id":   "pr-101",
			"pull_request_name": "Initial feature",
			"author_id":         "u1",
		}

		t.Run("1_CreatePR_Success", func(t *testing.T) {
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/create", prPayload)
			require.Equal(http.StatusCreated, res.StatusCode)

			err := json.Unmarshal([]byte(body), &prResponse)
			require.NoError(err)
			createdPR = prResponse.PR

			require.Equal("pr-101", createdPR.PullRequestID)
			require.Equal("u1", createdPR.AuthorID)
			require.Equal("OPEN", createdPR.Status)
			require.NotNil(createdPR.CreatedAt)
			require.Nil(createdPR.MergedAt)
			require.NotEmpty(createdPR.AssignedReviewers)
			require.NotContains(createdPR.AssignedReviewers, "u1")
			require.NotContains(createdPR.AssignedReviewers, "u3")
			require.LessOrEqual(len(createdPR.AssignedReviewers), 2)

			originalReviewer = createdPR.AssignedReviewers[0]
		})

		t.Run("2_CreatePR_Duplicate", func(t *testing.T) {
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/create", prPayload)
			require.Equal(http.StatusConflict, res.StatusCode)

			err := json.Unmarshal([]byte(body), &errResponse)
			require.NoError(err)
			require.Equal("PR_EXISTS", errResponse.Error.Code)
		})

		t.Run("3_CreatePR_AuthorNotFound", func(t *testing.T) {
			payload := map[string]string{
				"pull_request_id":   "pr-102",
				"pull_request_name": "Another feature",
				"author_id":         "u99",
			}
			res, _ := tu.MakeRequest(t, url, "POST", "/pullRequest/create", payload)
			require.Equal(http.StatusNotFound, res.StatusCode)
		})

		t.Run("4_GetReviews_Initial", func(t *testing.T) {
			path := fmt.Sprintf("/users/getReview?user_id=%s", originalReviewer)
			res, body := tu.MakeRequest(t, url, "GET", path, nil)
			require.Equal(http.StatusOK, res.StatusCode)

			err := json.Unmarshal([]byte(body), &reviewListResponse)
			require.NoError(err)
			require.Equal(originalReviewer, reviewListResponse.UserID)
			require.Len(reviewListResponse.PullRequests, 1)
			require.Equal("pr-101", reviewListResponse.PullRequests[0].PullRequestID)
			require.Equal("OPEN", reviewListResponse.PullRequests[0].Status)
		})

		t.Run("5_Reassign_Success", func(t *testing.T) {
			payload := map[string]string{
				"pull_request_id": "pr-101",
				"old_reviewer_id": originalReviewer,
			}
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/reassign", payload)
			require.Equal(http.StatusOK, res.StatusCode)

			err := json.Unmarshal([]byte(body), &reassignResponse)
			require.NoError(err)
			require.NotEmpty(reassignResponse.ReplacedBy)
			require.NotEqual(originalReviewer, reassignResponse.ReplacedBy)
			require.NotEqual("u1", reassignResponse.ReplacedBy)
			require.NotEqual("u3", reassignResponse.ReplacedBy)
			require.Contains(reassignResponse.PR.AssignedReviewers, reassignResponse.ReplacedBy)
			require.NotContains(reassignResponse.PR.AssignedReviewers, originalReviewer)

			createdPR = reassignResponse.PR
		})

		t.Run("6_Reassign_NotAssigned", func(t *testing.T) {
			payload := map[string]string{
				"pull_request_id": "pr-101",
				"old_reviewer_id": originalReviewer,
			}
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/reassign", payload)
			require.Equal(http.StatusConflict, res.StatusCode)

			err := json.Unmarshal([]byte(body), &errResponse)
			require.NoError(err)
			require.Equal("NOT_ASSIGNED", errResponse.Error.Code)
		})

		t.Run("7_Merge_Success", func(t *testing.T) {
			payload := map[string]string{"pull_request_id": "pr-101"}
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/merge", payload)
			require.Equal(http.StatusOK, res.StatusCode)

			err := json.Unmarshal([]byte(body), &prResponse)
			require.NoError(err)
			require.Equal("MERGED", prResponse.PR.Status)
			require.NotNil(prResponse.PR.MergedAt)

			createdPR = prResponse.PR
		})

		t.Run("8_Merge_Idempotent", func(t *testing.T) {
			payload := map[string]string{"pull_request_id": "pr-101"}
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/merge", payload)
			require.Equal(http.StatusOK, res.StatusCode)

			err := json.Unmarshal([]byte(body), &prResponse)
			require.NoError(err)
			require.Equal("MERGED", prResponse.PR.Status)
		})

		t.Run("9_Reassign_AfterMerge", func(t *testing.T) {
			newReviewer := createdPR.AssignedReviewers[0]
			payload := map[string]string{
				"pull_request_id": "pr-101",
				"old_reviewer_id": newReviewer,
			}
			res, body := tu.MakeRequest(t, url, "POST", "/pullRequest/reassign", payload)
			require.Equal(http.StatusConflict, res.StatusCode)

			err := json.Unmarshal([]byte(body), &errResponse)
			require.NoError(err)
			require.Equal("PR_MERGED", errResponse.Error.Code)
		})

		t.Run("10_Merge_NotFound", func(t *testing.T) {
			payload := map[string]string{"pull_request_id": "pr-999"}
			res, _ := tu.MakeRequest(t, url, "POST", "/pullRequest/merge", payload)
			require.Equal(http.StatusNotFound, res.StatusCode)
		})
	})
}
