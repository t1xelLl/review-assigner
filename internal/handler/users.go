package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type setIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) setIsActive(c *gin.Context) {
	var req setIsActiveRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	user, err := h.services.Teams.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{"user": user})
}

func (h *Handler) getReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		errorResp(c, models.NewAppError("INVALID_INPUT", "user_id is required"))
		return
	}

	prs, err := h.services.Users.GetUserReviewPRs(userID)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
