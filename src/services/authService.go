package services

import (
	"gin_main/pkg/hash"
	"gin_main/src/jwt"
	"gin_main/src/models"
	"net/http"

	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	ValidateBearerToken(token string) (models.User, error)
	LoginByBearerToken(user models.User) (string, *models.ErrorResponse)
	Registration(user models.RegistrationUserRequest) *models.ErrorResponse
}

type authService struct {
	jwtHelper   *jwt.JWTHelper
	userService UserServiceInterface
}

func NewAuthService(jwtHelper *jwt.JWTHelper, userService UserServiceInterface) AuthServiceInterface {
	return &authService{jwtHelper: jwtHelper, userService: userService}
}

func (s *authService) ValidateBearerToken(token string) (models.User, error) {
	var err error
	claims, err := s.jwtHelper.ValidateJWT(token)
	if err != nil {
		return models.User{}, err
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: userId, Login: claims.Login}, nil
}

func (s *authService) LoginByBearerToken(user models.User) (string, *models.ErrorResponse) {
	var err error
	userFromStorage, inError := s.userService.FindByLogin(user.Login)
	if inError != nil {
		return "", inError
	}
	if err = hash.VerifyPassword(userFromStorage.HashedPassword, user.Password); err != nil {
		return "", &models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	}
	token, err := s.jwtHelper.GenerateJWT(user)
	if err != nil {
		return "", &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return token, nil
}

func (s *authService) Registration(userData models.RegistrationUserRequest) *models.ErrorResponse {
	user := models.User{
		Login:    userData.Login,
		Password: userData.Password,
	}
	return s.userService.Create(user)
}
