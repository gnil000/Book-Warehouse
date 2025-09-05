package router

import (
	"gin_main/pkg/httpserver/middlewares"
	"gin_main/src/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type HandlerInterface interface {
	RegisterRoutes(api *gin.RouterGroup)
}

func RegisterPublicEndpoints(router *gin.Engine, handlers ...HandlerInterface) {
	api := router.Group("/")
	for _, handler := range handlers {
		handler.RegisterRoutes(api)
	}
}

func RegisterProtectedEndpoints(router *gin.Engine, authService services.AuthServiceInterface, log zerolog.Logger, handlers ...HandlerInterface) {
	protected := router.Group("/")
	protected.Use(middlewares.BearerAuthMiddleware(authService, log))
	for _, handler := range handlers {
		handler.RegisterRoutes(protected)
	}
}
