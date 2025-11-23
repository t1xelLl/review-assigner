package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type createTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []models.TeamMember `json:"members"`
}

type deactivateTeamRequest struct {
	TeamName string `json:"team_name"`
}

func (h *Handler) createTeam(c *gin.Context) {
	var req createTeamRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	members := make([]models.User, len(req.Members))
	for i, member := range req.Members {
		members[i] = models.User{
			ID:       member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: req.TeamName,
		}
	}

	team := &models.Team{
		Name:    req.TeamName,
		Members: members,
	}

	if err := h.services.Teams.CreateTeam(team); err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusCreated, gin.H{"team": team})
}

func (h *Handler) getTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		errorResp(c, models.NewAppError("INVALID_INPUT", "team_name is required"))
		return
	}

	team, err := h.services.Teams.GetTeam(teamName)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, team)
}

func (h *Handler) deactivateTeam(c *gin.Context) {
	var req deactivateTeamRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	if err := h.services.Teams.DeactivateTeam(req.TeamName); err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{"message": "team deactivated"})
}

func (h *Handler) safeDeactivateTeam(c *gin.Context) {
	var req deactivateTeamRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	affectedPRs, err := h.services.Teams.SafeDeactivateTeam(req.TeamName)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	response := gin.H{
		"message":      "team safely deactivated",
		"affected_prs": affectedPRs,
	}

	successResponse(c, http.StatusOK, response)
}
