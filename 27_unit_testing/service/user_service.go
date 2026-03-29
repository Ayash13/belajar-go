package service

import (
	"context"

	"belajar-go/27_unit_testing/domain"
)

type UserService struct {
	repo     domain.UserRepository
	producer domain.EventPublisher
}

func NewUserService(repo domain.UserRepository, producer domain.EventPublisher) *UserService {
	return &UserService{
		repo:     repo,
		producer: producer,
	}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.producer.PublishEvent(ctx, "user-fetched")
	return user, nil
}
