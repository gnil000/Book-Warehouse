package handlers

import (
	"gin_main/pkg/httpserver/router"
	"gin_main/src/models"
	"gin_main/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
	router.HandlerInterface
	Registration(ctx *gin.Context)
	//Login(ctx *gin.Context)
}

type authHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{authService: authService}
}

func (h *authHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/authentication/registration", h.Registration) //.
	//POST("/authentication/login", h.Login)
}

// @Summary      Регистрация пользователя
// @Description  Регистрирует пользователя в системе
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        user  body      models.RegistrationUserRequest  true  "Данные для регистрации пользователя"
// @Success      200
// @Failure      400     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /authentication/registration [post]
func (h *authHandler) Registration(ctx *gin.Context) {
	var user models.RegistrationUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if inError := h.authService.Registration(user); inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.Status(http.StatusOK)
}

// @Summary      Вход в систему
// @Description  Вход пользователя в систему по логину и паролю. Возвращает jwt токен
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "Данные для входа пользователя"
// @Success      200 {string} string "JWT токен"
// @Failure      400     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /authentication/login [post]
func (h *authHandler) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, inError := h.authService.LoginByBearerToken(user)
	if inError != nil {
		ctx.AbortWithStatusJSON(inError.Code, inError)
		return
	}
	ctx.String(http.StatusOK, token)
}
