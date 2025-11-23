package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type Teams interface {
	CreateTeam(team *models.Team) error
	GetTeam(name string) (*models.Team, error)
	UpdateUserActivity(userID string, isActive bool) error
	GetTeamMembers(teamName string) ([]models.User, error)
	DeactivateTeamMembers(teamName string) error
	SafeDeactivateTeamMembers(teamName string) ([]string, error)
}

type Users interface {
	GetUser(userID string) (*models.User, error)
	GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error)
	GetUserReviewPRs(userID string) ([]models.PullRequestShort, error)
}

type PullRequests interface {
	CreatePR(pr *models.PullRequest, reviewers []string) error
	GetPR(prID string) (*models.PullRequest, error)
	UpdatePRStatus(prID string, status models.PullRequestStatus, mergedAt *time.Time) error
	GetPRReviewers(prID string) ([]string, error)
	UpdatePRReviewers(prID string, oldReviewer, newReviewer string) error
	GetPRStats() (map[string]int, error)
	GetReviewerStats() (map[string]int, error)
}

type Repository struct {
	Teams
	Users
	PullRequests
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Teams:        NewTeamRepo(db),
		Users:        NewUserRepo(db),
		PullRequests: NewPRRepo(db),
	}
}
