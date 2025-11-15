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

type UserHandler struct {
	userSvc usecases.UserService
	pathCfg config.PathConfig
}

func NewUserHandler(
	userSvc usecases.UserService,
	pathCfg config.PathConfig,
) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		pathCfg: pathCfg,
	}
}

func (h *UserHandler) WithUserHandlers() handlers.RouterOption {
	return func (r chi.Router) {
		r.Post(h.pathCfg.SetIsActiveUser, h.setIsActiveHandler)
		r.Get(h.pathCfg.GetReviewUser, h.getReviewHandler)
	}
}

func (h *UserHandler) setIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateSetIsActiveRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.userSvc.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateSetIsActiveResponse(res))
}

func (h *UserHandler) getReviewHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetReviewRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.userSvc.GetReview(r.Context(), req.UserID)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateGetReviewResponse(req.UserID, res))
}
