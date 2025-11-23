package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type createPRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type mergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type reassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (h *Handler) createPR(c *gin.Context) {
	var req createPRRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	pr, err := h.services.PullRequests.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusCreated, gin.H{"pr": pr})
}

func (h *Handler) mergePR(c *gin.Context) {
	var req mergePRRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	pr, err := h.services.PullRequests.MergePR(req.PullRequestID)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) reassignReviewer(c *gin.Context) {
	var req reassignReviewerRequest
	if err := c.BindJSON(&req); err != nil {
		errorResp(c, models.NewAppError("INVALID_INPUT", "invalid input data"))
		return
	}

	pr, newReviewer, err := h.services.PullRequests.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		errorResp(c, err.(models.AppError))
		return
	}

	successResponse(c, http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}
