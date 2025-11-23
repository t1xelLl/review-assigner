package service

import (
	"github.com/t1xelLl/review-assigner/internal/models"
	"github.com/t1xelLl/review-assigner/internal/repository"
)

type TeamService struct {
	teamRepo repository.Teams
	userRepo repository.Users
}

func NewTeamService(teamRepo repository.Teams, userRepo repository.Users) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(team *models.Team) error {
	return s.teamRepo.CreateTeam(team)
}

func (s *TeamService) GetTeam(name string) (*models.Team, error) {
	return s.teamRepo.GetTeam(name)
}

func (s *TeamService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	err := s.teamRepo.UpdateUserActivity(userID, isActive)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetUser(userID)
}

func (s *TeamService) DeactivateTeam(teamName string) error {
	return s.teamRepo.DeactivateTeamMembers(teamName)
}

func (s *TeamService) SafeDeactivateTeam(teamName string) ([]string, error) {
	affectedPRs, err := s.teamRepo.SafeDeactivateTeamMembers(teamName)
	return affectedPRs, err
}
