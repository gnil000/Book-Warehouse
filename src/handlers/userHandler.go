package handlers

import (
	"gin_main/pkg/httpserver/router"
	"gin_main/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	router.HandlerInterface
	FindByLogin(ctx *gin.Context)
}

type userHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) UserHandlerInterface {
	return &userHandler{userService: userService}
}

func (h *userHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/users/:login", h.FindByLogin)
}

func (h *userHandler) FindByLogin(ctx *gin.Context) {
	login := ctx.Param("login")
	user, inError := h.userService.FindByLogin(login)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.JSON(http.StatusOK, user)
}
