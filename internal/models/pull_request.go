package models

import (
	"time"
)

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                string            `db:"id" json:"pull_request_id"`
	Name              string            `db:"name" json:"pull_request_name"`
	AuthorID          string            `db:"author_id" json:"author_id"`
	Status            PullRequestStatus `db:"status" json:"status"`
	AssignedReviewers []string          `db:"-" json:"assigned_reviewers"`
	CreatedAt         time.Time         `db:"created_at" json:"createdAt,omitempty"`
	MergedAt          *time.Time        `db:"merged_at" json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	ID     string            `db:"pull_request_id" json:"pull_request_id"`
	Name   string            `db:"pull_request_name" json:"pull_request_name"`
	Author string            `db:"author_id" json:"author_id"`
	Status PullRequestStatus `db:"status" json:"status"`
}
