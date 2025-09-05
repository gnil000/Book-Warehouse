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
