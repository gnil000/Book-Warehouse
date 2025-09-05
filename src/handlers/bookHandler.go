package handlers

import (
	"gin_main/src/models"
	"gin_main/src/services"
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
	router.GET("/books", h.GetAllBooks).
		GET("/books/:id", h.FindBookById).
		GET("/books/search", h.FindBookByParameters).
		POST("/books", h.CreateBook).
		PUT("/books", h.UpdateBook).
		PUT("/books/count", h.ChangeQuantity)
}

// @Summary Создать книгу
// @Description Создаёт объект описания книги и сохраняет в хранилище
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.CreateBookRequest true "Данные для создания книги"
// @Success 200 {object} models.CreateBookResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *bookHandler) CreateBook(ctx *gin.Context) {
	var createBookRequest models.CreateBookRequest
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

// @Summary Обновить книгу
// @Description Обновляет описание книги и сохраняет в хранилище
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.UpdateBookRequest true "Данные для обновления книги"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [put]
func (h *bookHandler) UpdateBook(ctx *gin.Context) {
	var updateBookRequest models.UpdateBookRequest
	if err := ctx.ShouldBindJSON(&updateBookRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inError := h.bookService.Update(updateBookRequest)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary Поиск книги по id
// @Description Ищет книгу по id в хранилище
// @Tags books
// @Produce json
// @Param book_id path uuid true "Id книги"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{book_id} [get]
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

// @Summary Поиск книги по параметрам
// @Description Ищет книгу по переданным в query параметрам
// @Tags books
// @Produce json
// @Param title query string false "Название книги"
// @Param author query string false "Имя автора (ФИО или Фамилия)"
// @Param yearOfWriting query string false "Год написания книги (формат YYYY-MM-DD)"
// @Param yearOfBirth query string false "Год рождения автора (формат YYYY-MM-DD)"
// @Success 200 {array} models.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/search [get]
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

// @Summary Получить все книги
// @Description Возвращает список всех книг в хранилище
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Failure 500 {object} map[string]string
// @Router /books [get]
func (h *bookHandler) GetAllBooks(ctx *gin.Context) {
	books, err := h.bookService.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, err)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// @Summary Изменить количество
// @Description Меняет количество копий конкретной книги в хранилище
// @Tags books
// @Produce json
// @Param quantity body models.ChangeBookQuantityRequest true "Id книги и количество для изменения (может быть больше или меньше нуля)"
// @Success 200 {array} models.ChangeBookQuantityResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/count [put]
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
