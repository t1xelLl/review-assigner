package service

import (
	"math/rand"
	"sync"
	"time"

	"github.com/t1xelLl/review-assigner/internal/models"
	"github.com/t1xelLl/review-assigner/internal/repository"
)

type PRService struct {
	prRepo   repository.PullRequests
	teamRepo repository.Teams
	userRepo repository.Users
}

func NewPRService(prRepo repository.PullRequests, teamRepo repository.Teams, userRepo repository.Users) *PRService {
	return &PRService{
		prRepo:   prRepo,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *PRService) CreatePR(prID, prName, authorID string) (*models.PullRequest, error) {
	author, err := s.userRepo.GetUser(authorID)
	if err != nil {
		return nil, err
	}

	candidates, err := s.userRepo.GetActiveTeamMembers(author.TeamName, authorID)
	if err != nil {
		return nil, err
	}

	reviewers := selectRandomReviewers(candidates, 2)

	reviewerIDs := make([]string, len(reviewers))
	for i, reviewer := range reviewers {
		reviewerIDs[i] = reviewer.ID
	}

	pr := &models.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		Status:            models.PullRequestStatusOpen,
		AssignedReviewers: reviewerIDs,
		CreatedAt:         time.Now(),
	}

	err = s.prRepo.CreatePR(pr, reviewerIDs)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) MergePR(prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetPR(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == models.PullRequestStatusMerged {
		return pr, nil
	}

	now := time.Now()
	err = s.prRepo.UpdatePRStatus(prID, models.PullRequestStatusMerged, &now)
	if err != nil {
		return nil, err
	}

	pr.Status = models.PullRequestStatusMerged
	pr.MergedAt = &now
	return pr, nil
}

func (s *PRService) ReassignReviewer(prID, oldUserID string) (*models.PullRequest, string, error) {
	pr, err := s.prRepo.GetPR(prID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == models.PullRequestStatusMerged {
		return nil, "", models.ErrPRMerged
	}

	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", models.ErrNotAssigned
	}

	oldUser, err := s.userRepo.GetUser(oldUserID)
	if err != nil {
		return nil, "", err
	}

	candidates, err := s.userRepo.GetActiveTeamMembers(oldUser.TeamName, oldUserID)
	if err != nil {
		return nil, "", err
	}

	filteredCandidates := make([]models.User, 0)
	for _, candidate := range candidates {
		if candidate.ID != pr.AuthorID {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	if len(filteredCandidates) == 0 {
		return nil, "", models.ErrNoCandidate
	}

	newReviewer := selectRandomReviewers(filteredCandidates, 1)[0]

	err = s.prRepo.UpdatePRReviewers(prID, oldUserID, newReviewer.ID)
	if err != nil {
		return nil, "", err
	}

	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			pr.AssignedReviewers[i] = newReviewer.ID
			break
		}
	}

	return pr, newReviewer.ID, nil
}

func (s *PRService) GetPRStats() (map[string]int, error) {
	return s.prRepo.GetPRStats()
}

func (s *PRService) GetReviewerStats() (map[string]int, error) {
	return s.prRepo.GetReviewerStats()
}

var (
	rng      = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMutex sync.Mutex
)

func selectRandomReviewers(candidates []models.User, max int) []models.User {
	if len(candidates) == 0 {
		return []models.User{}
	}

	if len(candidates) <= max {
		return candidates
	}

	rngMutex.Lock()
	defer rngMutex.Unlock()

	shuffled := make([]models.User, len(candidates))
	copy(shuffled, candidates)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:max]
}
