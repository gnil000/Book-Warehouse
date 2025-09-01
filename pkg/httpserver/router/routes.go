package router

import (
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
