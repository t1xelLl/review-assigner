package service

import (
	"github.com/t1xelLl/review-assigner/internal/models"
	"github.com/t1xelLl/review-assigner/internal/repository"
)

type Teams interface {
	CreateTeam(team *models.Team) error
	GetTeam(name string) (*models.Team, error)
	SetUserActive(userID string, isActive bool) (*models.User, error)
	DeactivateTeam(teamName string) error
	SafeDeactivateTeam(teamName string) ([]string, error)
}

type Users interface {
	GetUserReviewPRs(userID string) ([]models.PullRequestShort, error)
}

type PullRequests interface {
	CreatePR(prID, prName, authorID string) (*models.PullRequest, error)
	MergePR(prID string) (*models.PullRequest, error)
	ReassignReviewer(prID, oldUserID string) (*models.PullRequest, string, error)
	GetPRStats() (map[string]int, error)
	GetReviewerStats() (map[string]int, error)
}

type Service struct {
	Teams
	Users
	PullRequests
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Teams:        NewTeamService(repos.Teams, repos.Users),
		Users:        NewUserService(repos.Users),
		PullRequests: NewPRService(repos.PullRequests, repos.Teams, repos.Users),
	}
}
