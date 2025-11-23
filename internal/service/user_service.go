package service

import (
	"github.com/t1xelLl/review-assigner/internal/models"
	"github.com/t1xelLl/review-assigner/internal/repository"
)

type UserService struct {
	userRepo repository.Users
}

func NewUserService(userRepo repository.Users) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserReviewPRs(userID string) ([]models.PullRequestShort, error) {
	return s.userRepo.GetUserReviewPRs(userID)
}
