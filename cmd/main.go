package main

// https://gin-gonic.com/docs/

import (
	"fmt"
	"net/http"

	"gin_main/config"
	"gin_main/pkg/httpserver"
	"gin_main/pkg/httpserver/middlewares"
	"gin_main/pkg/httpserver/router"
	"gin_main/pkg/logger"
	"gin_main/src/database/migrations"
	"gin_main/src/database/repositories"
	"gin_main/src/handlers"
	"gin_main/src/jwt"
	"gin_main/src/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config := config.NewConfig()
	log := logger.NewLogger(zerolog.InfoLevel)

	migrations.Migrate(config, log)

	engine := gin.Default()
	server := httpserver.NewServer(log, engine, config)

	//	db := database.NewDatabaseConnection(config)
	//	bookRepo := repositories.NewBookRepository(db)
	//	bookService := services.NewBookService(bookRepo)
	//	bookHandler := handlers.NewBookHandler(bookService)

	//	authorRepo := repositories.NewAuthorRepository(db)
	//	authorService := services.NewAuthorService(authorRepo)
	//authorHandler := handlers.NewAuthorHandler(authorService)

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	//userHandler := handlers.NewUserHandler(userService)

	jwtHelper := jwt.NewJWTHelper(config, log)
	authService := services.NewAuthService(jwtHelper, userService, log)
	authHandler := handlers.NewAuthHandler(authService)

	server.AddMiddleware(middlewares.LogContextMiddleware(server.GetLogger()))
	//server.AddMiddleware(middlewares.BearerAuthMiddleware(authService))

	router.RegisterPublicEndpoints(engine, authHandler)
	// router.RegisterProtectedEndpoints(engine, authService, log, bookHandler)
	// router.RegisterProtectedEndpoints(engine, authService, log, authorHandler)
	// router.RegisterProtectedEndpoints(engine, authService, log, userHandler)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	engine.POST("/try", SomeHandler)
	server.AddSwagger()
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
