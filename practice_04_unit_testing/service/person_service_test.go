package service

import (
	"context"
	"testing"

	"practice_04_unit_testing/dto"
	"practice_04_unit_testing/entity"
	"practice_04_unit_testing/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersonService_CreatePerson(t *testing.T) {
	type Mocker struct {
		repo *mocks.PersonRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		input     dto.PersonCreateRequest
		wantErr   bool
		expected  dto.PersonResponse
	}{
		{
			desc: "SUCCESS: Create Person",
			input: dto.PersonCreateRequest{
				Name:  "ay",
				Email: "ay@mail.com",
			},
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"CreatePerson",
					mock.Anything,
					mock.AnythingOfType("*entity.Person"),
				).Run(func(args mock.Arguments) {
					person := args.Get(1).(*entity.Person)
					person.ID = 1
				}).Return(nil)
			},
			expected: dto.PersonResponse{
				ID:    1,
				Name:  "ay",
				Email: "ay@mail.com",
			},
		},
		{
			desc: "ERROR: Repository Returns Error",
			input: dto.PersonCreateRequest{
				Name:  "ay",
				Email: "ay@mail.com",
			},
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"CreatePerson",
					mock.Anything,
					mock.AnythingOfType("*entity.Person"),
				).Return(assert.AnError)
			},
			expected: dto.PersonResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := new(mocks.PersonRepository)
			mocker := &Mocker{repo: repo}

			tc.mockSetup(mocker)

			service := NewPersonService(repo)

			actual, err := service.CreatePerson(context.Background(), tc.input)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestPersonService_GetPerson(t *testing.T) {
	type Mocker struct {
		repo *mocks.PersonRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		inputID   int
		wantErr   bool
		expected  dto.PersonResponse
	}{
		{
			desc:    "SUCCESS: Get Person by ID",
			inputID: 1,
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetPerson",
					mock.Anything,
					mock.AnythingOfType("int"),
				).Return(&entity.Person{
					ID:    1,
					Name:  "ay",
					Email: "ay@mail.com",
				}, nil)
			},
			expected: dto.PersonResponse{
				ID:    1,
				Name:  "ay",
				Email: "ay@mail.com",
			},
		},
		{
			desc:    "ERROR: Person Not Found",
			inputID: 999,
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetPerson",
					mock.Anything,
					mock.AnythingOfType("int"),
				).Return(nil, assert.AnError)
			},
			expected: dto.PersonResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := new(mocks.PersonRepository)
			mocker := &Mocker{repo: repo}

			tc.mockSetup(mocker)

			service := NewPersonService(repo)

			actual, err := service.GetPerson(context.Background(), tc.inputID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestPersonService_GetAllPersons(t *testing.T) {
	type Mocker struct {
		repo *mocks.PersonRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  []dto.PersonResponse
	}{
		{
			desc:    "SUCCESS: Get All Persons",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetAllPersons",
					mock.Anything,
				).Return([]entity.Person{
					{ID: 1, Name: "ay", Email: "ay@mail.com"},
					{ID: 2, Name: "Jane Doe", Email: "jane@mail.com"},
				}, nil)
			},
			expected: []dto.PersonResponse{
				{ID: 1, Name: "ay", Email: "ay@mail.com"},
				{ID: 2, Name: "Jane Doe", Email: "jane@mail.com"},
			},
		},
		{
			desc:    "SUCCESS: Get All Persons Empty List",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetAllPersons",
					mock.Anything,
				).Return([]entity.Person{}, nil)
			},
			expected: nil,
		},
		{
			desc:    "ERROR: Repository Returns Error",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetAllPersons",
					mock.Anything,
				).Return(nil, assert.AnError)
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := new(mocks.PersonRepository)
			mocker := &Mocker{repo: repo}

			tc.mockSetup(mocker)

			service := NewPersonService(repo)

			actual, err := service.GetAllPersons(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			repo.AssertExpectations(t)
		})
	}
}
