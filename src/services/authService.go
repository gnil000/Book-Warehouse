package services

import (
	"fmt"
	"gin_main/pkg/hash"
	"gin_main/src/jwt"
	"gin_main/src/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type AuthServiceInterface interface {
	ValidateBearerToken(token string) (models.User, error)
	LoginByBearerToken(user models.User) (string, *models.ErrorResponse)
	Registration(user models.RegistrationUserRequest) *models.ErrorResponse
}

type authService struct {
	jwtHelper   *jwt.JWTHelper
	userService UserServiceInterface
	log         zerolog.Logger
}

func NewAuthService(jwtHelper *jwt.JWTHelper, userService UserServiceInterface, log zerolog.Logger) AuthServiceInterface {
	return &authService{jwtHelper: jwtHelper, userService: userService, log: log}
}

func (s *authService) ValidateBearerToken(token string) (models.User, error) {
	var err error
	claims, err := s.jwtHelper.ValidateJWT(token)
	s.log.Err(err)
	if err != nil {
		return models.User{}, err
	}
	s.log.Info().Msg(fmt.Sprintf("объект claims %v", *claims))
	userId, err := uuid.Parse(claims.Subject)
	s.log.Err(err)
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
