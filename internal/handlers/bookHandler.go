package handlers

import (
	"gin_main/internal/models"
	"gin_main/internal/services"
	"net/http"
	"time"

	"gin_main/pkg/httpserver/router"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookHandlerInterface interface {
	router.HandlerInterface
	CreateBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	FindBookById(ctx *gin.Context)
	FindBookByParameters(ctx *gin.Context)
	GetAllBooks(ctx *gin.Context)
	ChangeQuantity(ctx *gin.Context)
}

type bookHandler struct {
	bookService services.BookServiceInterface
}

func NewBookHandler(bookService services.BookServiceInterface) BookHandlerInterface {
	return &bookHandler{bookService: bookService}
}

func (h *bookHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/books", h.GetAllBooks)
}

func (h *bookHandler) CreateBook(ctx *gin.Context) {
	var createBookRequest models.CreateOrUpdateBookRequest
	if err := ctx.ShouldBindJSON(&createBookRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createBookResponse, inError := h.bookService.Create(createBookRequest)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, createBookResponse)
}

func (h *bookHandler) UpdateBook(ctx *gin.Context) {
	var updateBookRequest models.CreateOrUpdateBookRequest
	if err := ctx.ShouldBindJSON(&updateBookRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inError := h.bookService.Update(updateBookRequest)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *bookHandler) FindBookById(ctx *gin.Context) {
	var err error
	bookID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "book id not valid"})
		return
	}
	book, errResponse := h.bookService.FindById(bookID)
	if errResponse != nil {
		ctx.AbortWithStatusJSON(errResponse.Code, err)
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (h *bookHandler) FindBookByParameters(ctx *gin.Context) {
	title := ctx.Query("title")
	author := ctx.Query("author")
	var yearOfWritingPtr *time.Time = nil
	var yearOfBirthPtr *time.Time = nil
	if ctx.Query("yearOfWriting") != "" {
		yearOfWriting, err := time.Parse(time.DateOnly, ctx.Query("yearOfWriting"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot convert date yearOfWriting"})
			return
		}
		yearOfWritingPtr = &yearOfWriting
	}
	if ctx.Query("yearOfBirth") != "" {
		yearOfBirth, err := time.Parse(time.DateOnly, ctx.Query("yearOfBirth"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot convert date yearOfBirth"})
			return
		}
		yearOfBirthPtr = &yearOfBirth
	}
	books, inError := h.bookService.FindByParameters(title, author, yearOfWritingPtr, yearOfBirthPtr)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *bookHandler) GetAllBooks(ctx *gin.Context) {
	books, err := h.bookService.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, err)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *bookHandler) ChangeQuantity(ctx *gin.Context) {
	var changeQuantity models.ChangeBookQuantityRequest
	if err := ctx.ShouldBindJSON(&changeQuantity); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	changeQuantityResponse, inError := h.bookService.ChangeQuantity(changeQuantity)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, changeQuantityResponse)
}
