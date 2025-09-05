package router

import (
	"gin_main/pkg/httpserver/middlewares"
	"gin_main/src/services"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	RegisterRoutes(api *gin.RouterGroup)
}

func RegisterPublicEndpoints(router *gin.Engine, handlers ...HandlerInterface) {
	api := router.Group("/api")
	for _, handler := range handlers {
		handler.RegisterRoutes(api)
	}
}

func RegisterProtectedEndpoints(router *gin.Engine, authService services.AuthServiceInterface, handlers ...HandlerInterface) {
	protected := router.Group("/api")
	protected.Use(middlewares.BearerAuthMiddleware(authService))
	for _, handler := range handlers {
		handler.RegisterRoutes(protected)
	}
}
