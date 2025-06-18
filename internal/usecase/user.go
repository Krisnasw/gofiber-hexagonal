package usecase

import "app-hexagonal/internal/domain"

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
