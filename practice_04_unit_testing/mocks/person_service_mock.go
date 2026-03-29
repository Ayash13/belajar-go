package mocks

import (
	"context"
	"practice_04_unit_testing/dto"

	"github.com/stretchr/testify/mock"
)

type PersonService struct {
	mock.Mock
}

func (m *PersonService) CreatePerson(ctx context.Context, req dto.PersonCreateRequest) (dto.PersonResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.PersonResponse), args.Error(1)
}

func (m *PersonService) GetPerson(ctx context.Context, id int) (dto.PersonResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dto.PersonResponse), args.Error(1)
}

func (m *PersonService) GetAllPersons(ctx context.Context) ([]dto.PersonResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.PersonResponse), args.Error(1)
}
