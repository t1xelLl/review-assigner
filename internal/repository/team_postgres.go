package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type TeamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) CreateTeam(team *models.Team) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO teams (name) VALUES ($1)", team.Name)
	if err != nil {
		return models.ErrTeamExists
	}

	for _, member := range team.Members {
		_, err = tx.NamedExec(`
			INSERT INTO users (id, username, team_name, is_active) 
			VALUES (:id, :username, :team_name, :is_active)
			ON CONFLICT (id) DO UPDATE SET 
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, member)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TeamRepo) GetTeam(name string) (*models.Team, error) {
	var team models.Team
	err := r.db.Get(&team, "SELECT name FROM teams WHERE name = $1", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	var members []models.User
	err = r.db.Select(&members, "SELECT id, username, team_name, is_active FROM users WHERE team_name = $1", name)
	if err != nil {
		return nil, err
	}

	team.Members = members
	return &team, nil
}

func (r *TeamRepo) UpdateUserActivity(userID string, isActive bool) error {
	result, err := r.db.Exec("UPDATE users SET is_active = $1 WHERE id = $2", isActive, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *TeamRepo) GetTeamMembers(teamName string) ([]models.User, error) {
	var members []models.User
	err := r.db.Select(&members,
		"SELECT id, username, team_name, is_active FROM users WHERE team_name = $1", teamName)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *TeamRepo) DeactivateTeamMembers(teamName string) error {
	_, err := r.db.Exec("UPDATE users SET is_active = false WHERE team_name = $1", teamName)
	return err
}

func (r *TeamRepo) SafeDeactivateTeamMembers(teamName string) ([]string, error) {
	var affectedPRs []string
	query := `
        SELECT DISTINCT pr.id 
        FROM pull_requests pr
        JOIN pr_reviewers prr ON pr.id = prr.pr_id  
        JOIN users u ON prr.user_id = u.id
        WHERE u.team_name = $1 AND pr.status = 'OPEN' AND u.is_active = true
    `
	err := r.db.Select(&affectedPRs, query, teamName)
	if err != nil {
		return nil, err
	}

	for _, prID := range affectedPRs {
		if err := r.safeReassignDeactivatedReviewers(prID, teamName); err != nil {
			logrus.Warnf("Failed to reassign PR %s: %v", prID, err)
		}
	}

	_, err = r.db.Exec("UPDATE users SET is_active = false WHERE team_name = $1", teamName)
	if err != nil {
		return nil, err
	}

	return affectedPRs, nil
}

func (r *TeamRepo) safeReassignDeactivatedReviewers(prID, teamName string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var deactivatedReviewers []string
	query := `
        SELECT prr.user_id
        FROM pr_reviewers prr
        JOIN users u ON prr.user_id = u.id
        WHERE prr.pr_id = $1 AND u.team_name = $2 AND u.is_active = true
    `
	err = tx.Select(&deactivatedReviewers, query, prID, teamName)
	if err != nil {
		return err
	}

	for _, oldReviewer := range deactivatedReviewers {
		var newReviewer string
		candidateQuery := `
            SELECT u.id 
            FROM users u
            WHERE u.team_name = $1 
            AND u.is_active = true 
            AND u.id != $2
            AND u.id NOT IN (
                SELECT user_id FROM pr_reviewers WHERE pr_id = $3
            )
            AND u.id != (
                SELECT author_id FROM pull_requests WHERE id = $3
            )
            LIMIT 1
        `
		err = tx.Get(&newReviewer, candidateQuery, teamName, oldReviewer, prID)
		if err == nil {
			_, err = tx.Exec(
				"UPDATE pr_reviewers SET user_id = $1 WHERE pr_id = $2 AND user_id = $3",
				newReviewer, prID, oldReviewer)
			if err != nil {
				return err
			}
		} else {
			_, err = tx.Exec(
				"DELETE FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2",
				prID, oldReviewer)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
