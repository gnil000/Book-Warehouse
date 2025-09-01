package httpserver

import (
	"context"
	"errors"
	"fmt"
	"gin_main/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	logger *zerolog.Logger
	router *gin.Engine
	config *config.Config
}

func NewServer(logger *zerolog.Logger, router *gin.Engine, config *config.Config) *Server {
	return &Server{logger: logger, router: router, config: config}
}

func (s *Server) Serve() {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:           s.router.Handler(),
		ReadTimeout:       s.config.Server.ReadTimeout,
		WriteTimeout:      s.config.Server.WriteTimeout,
		ReadHeaderTimeout: s.config.Server.ReadHeaderTimeout,
		IdleTimeout:       s.config.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal().Err(err).Msg("listen")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //INF: запихивает ожидаемые перечисленные сигналы в канал quit
	<-quit
	s.logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatal().Err(err).Msg("Server shutdown")
	}

	<-ctx.Done()
	s.logger.Info().Msg("Server shutdown timeout of 5 seconds")
	s.logger.Info().Msg("Server exiting")
}
