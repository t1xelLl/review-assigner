package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/review-assigner/internal/models"
)

type PRRepo struct {
	db *sqlx.DB
}

func NewPRRepo(db *sqlx.DB) *PRRepo {
	return &PRRepo{db: db}
}

func (r *PRRepo) CreatePR(pr *models.PullRequest, reviewers []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
		INSERT INTO pull_requests (id, name, author_id, status, created_at)
		VALUES (:id, :name, :author_id, :status, :created_at)
	`, pr)
	if err != nil {
		return models.ErrPRExists
	}

	for _, reviewer := range reviewers {
		_, err = tx.Exec("INSERT INTO pr_reviewers (pr_id, user_id) VALUES ($1, $2)",
			pr.ID, reviewer)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PRRepo) GetPR(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.Get(&pr,
		"SELECT id, name, author_id, status, created_at, merged_at FROM pull_requests WHERE id = $1", prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	reviewers, err := r.GetPRReviewers(prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PRRepo) GetPRReviewers(prID string) ([]string, error) {
	var reviewers []string
	err := r.db.Select(&reviewers,
		"SELECT user_id FROM pr_reviewers WHERE pr_id = $1", prID)
	if err != nil {
		return nil, err
	}
	return reviewers, nil
}

func (r *PRRepo) UpdatePRStatus(prID string, status models.PullRequestStatus, mergedAt *time.Time) error {
	var result sql.Result
	var err error

	if status == models.PullRequestStatusMerged {
		result, err = r.db.Exec(
			"UPDATE pull_requests SET status = $1, merged_at = $2 WHERE id = $3",
			status, mergedAt, prID)
	} else {
		result, err = r.db.Exec(
			"UPDATE pull_requests SET status = $1, merged_at = NULL WHERE id = $2",
			status, prID)
	}

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

func (r *PRRepo) UpdatePRReviewers(prID string, oldReviewer, newReviewer string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"DELETE FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2",
		prID, oldReviewer)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrNotAssigned
	}

	_, err = tx.Exec(
		"INSERT INTO pr_reviewers (pr_id, user_id) VALUES ($1, $2)",
		prID, newReviewer)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PRRepo) GetPRStats() (map[string]int, error) {
	stats := make(map[string]int)

	var statusStats []struct {
		Status string `db:"status"`
		Count  int    `db:"count"`
	}

	err := r.db.Select(&statusStats,
		"SELECT status, COUNT(*) as count FROM pull_requests GROUP BY status")
	if err != nil {
		return nil, err
	}

	for _, stat := range statusStats {
		stats[stat.Status] = stat.Count
	}

	return stats, nil
}

func (r *PRRepo) GetReviewerStats() (map[string]int, error) {
	stats := make(map[string]int)

	var reviewerStats []struct {
		UserID string `db:"user_id"`
		Count  int    `db:"count"`
	}

	err := r.db.Select(&reviewerStats,
		"SELECT user_id, COUNT(*) as count FROM pr_reviewers GROUP BY user_id")
	if err != nil {
		return nil, err
	}

	for _, stat := range reviewerStats {
		stats[stat.UserID] = stat.Count
	}

	return stats, nil
}
