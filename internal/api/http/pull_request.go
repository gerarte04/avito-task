package http

import (
	"avito-task/internal/api/http/response"
	"avito-task/internal/api/http/types"
	"avito-task/internal/config"
	"avito-task/internal/usecases"
	"avito-task/pkg/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PullRequestHandler struct {
	prSvc usecases.PullRequestService
	pathCfg config.PathConfig
}

func NewPullRequestHandler(
	prSvc usecases.PullRequestService,
	pathCfg config.PathConfig,
) *PullRequestHandler {
	return &PullRequestHandler{
		prSvc: prSvc,
		pathCfg: pathCfg,
	}
}

func (h *PullRequestHandler) WithPRHandlers() handlers.RouterOption {
	return func (r chi.Router) {
		r.Post(h.pathCfg.CreatePR, h.createHandler)
		r.Post(h.pathCfg.MergePR, h.mergeHandler)
		r.Post(h.pathCfg.ReassignPR, h.reassignHandler)
	}
}

func (h *PullRequestHandler) createHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.MakeCreatePRRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.prSvc.CreatePullRequest(r.Context(), req.PR)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusCreated, types.MakeCreatePRResponse(res))
}

func (h *PullRequestHandler) mergeHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateMergePRRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.prSvc.Merge(r.Context(), req.PRID)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateMergePRResponse(res))
}

func (h *PullRequestHandler) reassignHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateReassignRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	newRewID, pr, err := h.prSvc.Reassign(r.Context(), req.PRID, req.OldRewID)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateReassignResponse(newRewID, pr))
}
