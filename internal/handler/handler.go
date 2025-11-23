package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/review-assigner/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	team := router.Group("/team")
	{
		team.POST("/add", h.createTeam)
		team.GET("/get", h.getTeam)
		team.POST("/deactivate", h.deactivateTeam)
		team.POST("/safe-deactivate", h.safeDeactivateTeam)
	}

	users := router.Group("/users")
	{
		users.POST("/setIsActive", h.setIsActive)
		users.GET("/getReview", h.getReview)
	}

	pullRequests := router.Group("/pullRequests")
	{
		pullRequests.POST("/create", h.createPR)
		pullRequests.POST("/merge", h.mergePR)
		pullRequests.POST("/reassign", h.reassignReviewer)
	}

	stats := router.Group("/stats")
	{
		stats.GET("/reviewers", h.getReviewersStats)
		stats.GET("/pr", h.getPRStats)
	}

	return router
}
