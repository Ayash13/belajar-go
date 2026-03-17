package service

import (
	"context"
	"nethttp/dto"
	"nethttp/entity"
	"nethttp/repository"
)

type PersonService interface {
	CreatePerson(ctx context.Context, req dto.PersonCreateRequest) (dto.PersonResponse, error)
	GetPerson(ctx context.Context, id int) (dto.PersonResponse, error)
	GetAllPersons(ctx context.Context) ([]dto.PersonResponse, error)
}

type personServiceImpl struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) PersonService {
	return &personServiceImpl{repo: repo}
}

func (s *personServiceImpl) CreatePerson(ctx context.Context, req dto.PersonCreateRequest) (dto.PersonResponse, error) {
	person := &entity.Person{
		Name:  req.Name,
		Email: req.Email,
	}

	err := s.repo.CreatePerson(ctx, person)
	if err != nil {
		return dto.PersonResponse{}, err
	}

	return dto.PersonResponse{
		ID:    person.ID,
		Name:  person.Name,
		Email: person.Email,
	}, nil
}

func (s *personServiceImpl) GetPerson(ctx context.Context, id int) (dto.PersonResponse, error) {
	person, err := s.repo.GetPerson(ctx, id)
	if err != nil {
		return dto.PersonResponse{}, err
	}

	return dto.PersonResponse{
		ID:    person.ID,
		Name:  person.Name,
		Email: person.Email,
	}, nil
}

func (s *personServiceImpl) GetAllPersons(ctx context.Context) ([]dto.PersonResponse, error) {
	persons, err := s.repo.GetAllPersons(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, p := range persons {
		responses = append(responses, dto.PersonResponse{
			ID:    p.ID,
			Name:  p.Name,
			Email: p.Email,
		})
	}

	return responses, nil
}
