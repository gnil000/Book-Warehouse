package main

// https://gin-gonic.com/docs/

import (
	"fmt"
	"net/http"

	"gin_main/config"
	"gin_main/internal/handlers"
	"gin_main/internal/repositories"
	"gin_main/internal/services"
	"gin_main/pkg/database"
	"gin_main/pkg/httpserver"
	"gin_main/pkg/httpserver/router"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config := config.NewConfig()

	engine := gin.Default()

	server := httpserver.NewServer(&log.Logger, engine, config)

	db := database.NewDatabaseConnection(config)
	bookRepo := repositories.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	router.RegisterPublicEndpoints(engine, bookHandler)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	engine.POST("/try", SomeHandler)
	server.Serve()
}

func SomeHandler(c *gin.Context) {
	objA := formA{}
	objB := formB{}
	if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
		c.String(http.StatusOK, `the body should be formA`)
	}
	if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
		c.String(http.StatusOK, `the body should be formB JSON`)
	}
	fmt.Println("Return not was be called")
	fmt.Printf("data received: %+v, %+v", objA, objB)
}

type formA struct {
	Foo string `json:"foo" binding:"required"`
}

type formB struct {
	Bar string `json:"bar" binding:"required"`
}
