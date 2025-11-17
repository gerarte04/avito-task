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

type TeamHandler struct {
	teamSvc usecases.TeamService
	pathCfg config.PathConfig
}

func NewTeamHandler(
	teamSvc usecases.TeamService,
	pathCfg config.PathConfig,
) *TeamHandler {
	return &TeamHandler{
		teamSvc: teamSvc,
		pathCfg: pathCfg,
	}
}

func (h *TeamHandler) WithTeamHandlers() handlers.RouterOption {
	return func (r chi.Router) {
		r.Post(h.pathCfg.AddTeam, h.addHandler)
		r.Get(h.pathCfg.GetTeam, h.getHandler)
		r.Get(h.pathCfg.GetTeamStats, h.getStatsHandler)
		r.Post(h.pathCfg.DeactivateTeam, h.deactivateHandler)
	}
}

func (h *TeamHandler) addHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateAddTeamRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.teamSvc.CreateTeam(r.Context(), req.Team)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusCreated, types.CreateAddTeamResponse(res))
}

func (h *TeamHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTeamRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.teamSvc.GetTeam(r.Context(), req.Name)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateGetTeamResponse(res))
}

func (h *TeamHandler) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTeamStatsRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.teamSvc.GetTeamStats(r.Context(), req.Name)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, res)
}

func (h *TeamHandler) deactivateHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateDeactivateTeamRequest(r)
	if err != nil {
		response.ProcessCreatingRequestError(w, err)
		return
	}

	res, err := h.teamSvc.DeactivateTeam(r.Context(), req.Name)
	if err != nil {
		response.ProcessError(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, types.CreateDeactivateTeamResponse(req.Name, res))
}
