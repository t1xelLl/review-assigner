package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/models"
)

func (h *Handler) getReviewersStats(c *gin.Context) {
	stats, err := h.services.PullRequests.GetReviewerStats()
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{"reviewers": stats})
}

func (h *Handler) getPRStats(c *gin.Context) {
	stats, err := h.services.PullRequests.GetPRStats()
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{"pull_requests": stats})
}
