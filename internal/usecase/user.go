package usecase

import "app-hexagonal/internal/domain"

// UserUsecaseInterface defines the interface for user use cases
// This helps with dependency inversion in our hexagonal architecture
type UserUsecaseInterface interface {
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id string) error
}

type UserUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) GetUserByID(id string) (*domain.User, error) {
	return uc.repo.FindByID(id)
}

func (uc *UserUsecase) GetUserByEmail(email string) (*domain.User, error) {
	return uc.repo.FindByEmail(email)
}

func (uc *UserUsecase) CreateUser(user *domain.User) error {
	return uc.repo.Store(user)
}

func (uc *UserUsecase) UpdateUser(user *domain.User) error {
	return uc.repo.Update(user)
}

func (uc *UserUsecase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}
