package services

import (
	"fmt"
	"gin_main/src/database/entities"
	"gin_main/src/database/repositories"
	"gin_main/src/models"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AuthorServiceInterface interface {
	Create(author models.CreateAuthorRequest) (models.CreateAuthorResponse, *models.ErrorResponse)
	FindById(id uuid.UUID) (models.Author, *models.ErrorResponse)
	FindBySurname(surname string) ([]models.Author, *models.ErrorResponse)
	Update(author models.UpdateAuthorRequest) *models.ErrorResponse
}

type authorService struct {
	authorRepo repositories.AuthorRepositoryInterface
}

func NewAuthorService(authorRepo repositories.AuthorRepositoryInterface) AuthorServiceInterface {
	return &authorService{authorRepo: authorRepo}
}

func (s *authorService) Create(author models.CreateAuthorRequest) (models.CreateAuthorResponse, *models.ErrorResponse) {
	var err error
	var authorEntity entities.Author
	if err = copier.Copy(&authorEntity, &author); err != nil {
		return models.CreateAuthorResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = s.authorRepo.Create(&authorEntity); err != nil {
		return models.CreateAuthorResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	var authorResponse models.CreateAuthorResponse
	if err = copier.Copy(&authorResponse, &authorEntity); err != nil {
		return models.CreateAuthorResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return authorResponse, nil
}

func (s *authorService) FindById(id uuid.UUID) (models.Author, *models.ErrorResponse) {
	var err error
	var authorModel models.Author
	var authorEntity entities.Author
	if authorEntity, err = s.authorRepo.FindById(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.Author{}, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("author not found with id = %s", id.String()),
			}
		}
		return models.Author{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = copier.Copy(&authorModel, &authorEntity); err != nil {
		return models.Author{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return authorModel, nil
}

func (s *authorService) FindBySurname(surname string) ([]models.Author, *models.ErrorResponse) {
	var err error
	var authorModels []models.Author
	var authorEntities []entities.Author
	if authorEntities, err = s.authorRepo.FindBySurname(surname); err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = copier.Copy(&authorModels, &authorEntities); err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return authorModels, nil
}

func (s *authorService) Update(author models.UpdateAuthorRequest) *models.ErrorResponse {
	var err error
	var authorEntity entities.Author
	if err = copier.Copy(&authorEntity, &author); err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = s.authorRepo.Update(authorEntity); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("author not found with id = %s", author.ID.String()),
			}
		}
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
