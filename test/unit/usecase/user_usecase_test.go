package usecase_test

import (
	"errors"
	"testing"

	"app-hexagonal/internal/domain"
	"app-hexagonal/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) Store(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserUsecase_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:    "1",
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("FindByID", "1").Return(expectedUser, nil)

		user, err := userUsecase.GetUserByID("1")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("FindByID", "999").Return((*domain.User)(nil), errors.New("user not found"))

		user, err := userUsecase.GetUserByID("999")

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
