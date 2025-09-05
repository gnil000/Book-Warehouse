package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GetLogger(c *gin.Context) zerolog.Logger {
	loggerIface, exists := c.Get("logger")
	if !exists {
		panic("logger not found in context")
	}
	return loggerIface.(zerolog.Logger)
}
