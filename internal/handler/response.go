package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type errorResponse struct {
	Error models.AppError `json:"error"`
}

func successResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func errorResp(c *gin.Context, appErr models.AppError) {
	c.JSON(getHTTPCode(appErr.Code), errorResponse{Error: appErr})
}

func getHTTPCode(errorCode models.ErrorCode) int {
	switch errorCode {
	case models.ErrorCodeTeamExists, models.ErrorCodePullRExists:
		return http.StatusBadRequest
	case models.ErrorCodePRMerged, models.ErrorCodeNotAssigned, models.ErrorCodeNoCandidate:
		return http.StatusConflict
	case models.ErrorCodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
