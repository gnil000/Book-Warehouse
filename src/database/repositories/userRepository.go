package repositories

import (
	"fmt"
	"gin_main/src/database/entities"
)

type UserRepositoryInterface interface {
	Create(user entities.User) error
	FindByLogin(login string) (entities.User, error)
}

type userRepository struct {
	users []entities.User
}

func NewUserRepository() UserRepositoryInterface {
	return &userRepository{users: make([]entities.User, 0)}
}

func (r *userRepository) Create(user entities.User) error {
	r.users = append(r.users, user)
	return nil
}

func (r *userRepository) FindByLogin(login string) (entities.User, error) {
	for _, foundUser := range r.users {
		if foundUser.Login == login {
			return foundUser, nil
		}
	}
	return entities.User{}, fmt.Errorf("user with login %s not found", login)
}
