package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user,
		"SELECT id, username, team_name, is_active FROM users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error) {
	var users []models.User
	query := `
		SELECT id, username, team_name, is_active 
		FROM users 
		WHERE team_name = $1 AND is_active = true AND id != $2
		ORDER BY id
	`
	err := r.db.Select(&users, query, teamName, excludeUserID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetUserReviewPRs(userID string) ([]models.PullRequestShort, error) {
	var prs []models.PullRequestShort
	query := `
		SELECT p.id as pull_request_id, p.name as pull_request_name, 
		       p.author_id, p.status
		FROM pull_requests p
		JOIN pr_reviewers pr ON p.id = pr.pr_id
		WHERE pr.user_id = $1
		ORDER BY p.created_at DESC
	`
	err := r.db.Select(&prs, query, userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
