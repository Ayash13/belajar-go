package service

import (
	"context"
	"testing"

	"belajar-go/27_unit_testing/domain"
	"belajar-go/27_unit_testing/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetByID(t *testing.T) {
	type Mocker struct {
		repo         *mocks.UserRepository
		mockProducer *mocks.EventPublisher
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  *domain.User
	}{
		{
			desc:    "SUCCESS: Get User by ID",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				// Mock the repository call
				m.repo.On(
					"GetByID",
					mock.Anything,
					mock.AnythingOfType("int64"),
				).Return(&domain.User{
					ID:    100,
					Name:  "John Doe",
					Email: "dummy@mail.co",
				}, nil)

				// Mock the publisher call
				m.mockProducer.On(
					"PublishEvent",
					mock.Anything,
					"user-fetched",
				).Return(nil)
			},
			expected: &domain.User{
				ID:    100,
				Name:  "John Doe",
				Email: "dummy@mail.co",
			},
		},
		{
			desc:    "ERROR: User Not Found",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetByID",
					mock.Anything,
					mock.AnythingOfType("int64"),
				).Return(nil, assert.AnError)
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// 1. Initialize fresh mocks per test case
			repo := new(mocks.UserRepository)
			producer := new(mocks.EventPublisher)

			mocker := &Mocker{
				repo:         repo,
				mockProducer: producer,
			}

			// 2. Setup mock expectations
			tc.mockSetup(mocker)

			// 3. Initialize service with mocks
			service := NewUserService(repo, producer)

			// 4. Execute method
			actual, err := service.GetByID(context.Background(), 100)

			// 5. Assert results
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			// 6. Verify mock expectations
			repo.AssertExpectations(t)
			producer.AssertExpectations(t)
		})
	}
}
