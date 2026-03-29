package mocks

import (
	"context"
	"practice_04_unit_testing/entity"

	"github.com/stretchr/testify/mock"
)

type PersonRepository struct {
	mock.Mock
}

func (m *PersonRepository) CreatePerson(ctx context.Context, person *entity.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

func (m *PersonRepository) GetPerson(ctx context.Context, id int) (*entity.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Person), args.Error(1)
}

func (m *PersonRepository) GetAllPersons(ctx context.Context) ([]entity.Person, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Person), args.Error(1)
}
