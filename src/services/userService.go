package services

import (
	"gin_main/pkg/hash"
	"gin_main/src/database/entities"
	"gin_main/src/database/repositories"
	"gin_main/src/models"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
)

type UserServiceInterface interface {
	Create(user models.User) *models.ErrorResponse
	FindByLogin(login string) (models.User, *models.ErrorResponse)
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepository: userRepo}
}

func (s *userService) Create(user models.User) *models.ErrorResponse {
	var userEntity entities.User
	var err error
	if err = copier.Copy(&userEntity, &user); err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	userEntity.HashedPassword, err = hash.HashPassword(user.Password)
	if err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = s.userRepository.Create(userEntity); err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}

func (s *userService) FindByLogin(login string) (models.User, *models.ErrorResponse) {
	var userModel models.User
	userEntity, err := s.userRepository.FindByLogin(login)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.User{}, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "user not found",
			}
		}
		return models.User{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = copier.Copy(&userModel, &userEntity); err != nil {
		return models.User{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return userModel, nil
}
