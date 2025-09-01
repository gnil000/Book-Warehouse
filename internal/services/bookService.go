package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gin_main/internal/models"
	"gin_main/internal/repositories"
	"gin_main/internal/repositories/entities"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type BookServiceInterface interface {
	Create(book models.CreateOrUpdateBookRequest) (models.CreateBookResponse, *models.ErrorResponse)
	Update(book models.CreateOrUpdateBookRequest) *models.ErrorResponse
	FindById(id uuid.UUID) (models.Book, *models.ErrorResponse)
	FindByParameters(title, author string, yearOfWriting, yearOfBirth *time.Time) ([]models.Book, *models.ErrorResponse)
	GetAll() ([]models.Book, *models.ErrorResponse)
	ChangeQuantity(book models.ChangeBookQuantityRequest) (models.ChangeBookQuantityResponse, *models.ErrorResponse)
}

type bookService struct {
	bookRepo repositories.BookRepositoryInterface
}

func NewBookService(bookRepo repositories.BookRepositoryInterface) BookServiceInterface {
	return &bookService{bookRepo: bookRepo}
}

func (r *bookService) Create(book models.CreateOrUpdateBookRequest) (models.CreateBookResponse, *models.ErrorResponse) {
	var err error
	var bookEntity entities.Book
	if err = copier.Copy(&bookEntity, &book); err != nil {
		return models.CreateBookResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	newBookEntity, err := r.bookRepo.Create(bookEntity)
	if err != nil {
		return models.CreateBookResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	var bookResponse models.CreateBookResponse
	if err = copier.Copy(&bookResponse, &newBookEntity); err != nil {
		return models.CreateBookResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return bookResponse, nil
}

func (r *bookService) Update(book models.CreateOrUpdateBookRequest) *models.ErrorResponse {
	var err error
	var bookEntity entities.Book
	if err = copier.Copy(&bookEntity, &book); err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	if err := r.bookRepo.Update(bookEntity); err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return nil
}

func (r *bookService) FindById(id uuid.UUID) (models.Book, *models.ErrorResponse) {
	var err error
	bookFound, err := r.bookRepo.FindById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Book{}, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("user with id = %s not found", id.String()),
			}
		}
		return models.Book{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	var bookResult models.Book
	if err = copier.Copy(&bookResult, &bookFound); err != nil {
		return models.Book{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return bookResult, nil
}

func (r *bookService) FindByParameters(title, author string, yearOfWriting, yearOfBirth *time.Time) ([]models.Book, *models.ErrorResponse) {
	books, err := r.bookRepo.FindByParameters(title, author, yearOfWriting, yearOfBirth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Book{}, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "no one user found",
			}
		}
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	var booksResult []models.Book
	if err = copier.Copy(&booksResult, &books); err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return booksResult, nil
}

func (r *bookService) GetAll() ([]models.Book, *models.ErrorResponse) {
	var err error
	booksEntities, err := r.bookRepo.GetAll()
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	var books []models.Book
	if err = copier.Copy(&books, &booksEntities); err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return books, nil
}

func (r *bookService) ChangeQuantity(book models.ChangeBookQuantityRequest) (models.ChangeBookQuantityResponse, *models.ErrorResponse) {
	var err error
	newQuantity, err := r.bookRepo.ChangeQuantity(book.ID, book.Quantity)
	if err != nil {
		if strings.Contains(err.Error(), "negative") {
			return models.ChangeBookQuantityResponse{}, &models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "quantity cannot be negative",
			}
		}
		return models.ChangeBookQuantityResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return models.ChangeBookQuantityResponse{Quantity: newQuantity}, nil
}
