package handlers

import (
	"gin_main/pkg/httpserver/router"
	"gin_main/src/models"
	"gin_main/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorHandlerInterface interface {
	router.HandlerInterface
	CreateAuthor(ctx *gin.Context)
	FindAuthorById(ctx *gin.Context)
	FindAuthorBySurname(ctx *gin.Context)
	UpdateAuthor(ctx *gin.Context)
}

type authorHandler struct {
	authorService services.AuthorServiceInterface
}

func NewAuthorHandler(authorService services.AuthorServiceInterface) AuthorHandlerInterface {
	return &authorHandler{authorService: authorService}
}

func (h *authorHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/authors/:id", h.FindAuthorById).
		GET("/authors/search/:surname", h.FindAuthorBySurname).
		POST("/authors", h.CreateAuthor).
		PUT("/authors", h.UpdateAuthor)
}

// @Summary      Создать автора
// @Description  Добавляет нового автора в базу данных
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body      models.CreateAuthorRequest  true  "Данные для создания автора"
// @Success      200     {object}  models.CreateAuthorResponse
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /authors [post]
func (h *authorHandler) CreateAuthor(ctx *gin.Context) {
	var createAuthorRequest models.CreateAuthorRequest
	if err := ctx.ShouldBindJSON(&createAuthorRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createAuthorResponse, inError := h.authorService.Create(createAuthorRequest)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, createAuthorResponse)
}

// @Summary      Найти автора по ID
// @Description  Возвращает информацию об авторе по его UUID
// @Tags         authors
// @Produce      json
// @Param        author_id   path      string  true  "ID автора (UUID)"
// @Success      200  {object}  models.Author
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /authors/{id} [get]
func (h *authorHandler) FindAuthorById(ctx *gin.Context) {
	var err error
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id not valid"})
		return
	}
	author, inError := h.authorService.FindById(id)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, author)
}

// @Summary      Найти автора по фамилии
// @Description  Возвращает список авторов, у которых фамилия соответствует запросу
// @Tags         authors
// @Produce      json
// @Param        surname  path      string  true  "Фамилия автора"
// @Success      200      {array}   models.Author
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /authors/search/{surname} [get]
func (h *authorHandler) FindAuthorBySurname(ctx *gin.Context) {
	surname := ctx.Param("surname")
	// if surname == "" {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "path parameter surname is empty"})
	// 	return
	// }
	authors, inError := h.authorService.FindBySurname(surname)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, authors)
}

// @Summary      Обновить автора
// @Description  Обновляет информацию об авторе в базе данных
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body      models.UpdateAuthorRequest  true  "Данные для обновления автора"
// @Success      200
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /authors [put]
func (h *authorHandler) UpdateAuthor(ctx *gin.Context) {
	var updateAuthorRequest models.UpdateAuthorRequest
	if err := ctx.ShouldBindJSON(&updateAuthorRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inError := h.authorService.Update(updateAuthorRequest)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.Status(http.StatusOK)
}
