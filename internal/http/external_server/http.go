package external_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/rs/zerolog"
)

type IServer interface {
	Start() error
	Stop() error
}

type Server struct {
	config *config.ExternalServer
	logger *zerolog.Logger
	cache  cache.ICache
	router *gin.Engine
	http   http.Server
}

func NewServer(cfg *config.ExternalServer, logger *zerolog.Logger, cache cache.ICache) IServer {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.CustomRecoveryWithWriter(nil, recoveryFunc(logger)),
	)

	s := &Server{
		logger: logger,
		router: router,
		cache:  cache,
		http: http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: router,
		},
	}

	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	internalAuthGroup := s.router.Group("/auth/api/v1/internal")
	internalAuthGroup.Use(s.internalAuthMiddleWare)
	internalAuthGroup.POST("/new/oauth")

	externalAuthGroup := s.router.Group("/auth/api/v1/external")
	externalAuthGroup.POST("/new/sid")
}

func recoveryFunc(logger *zerolog.Logger) gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		logger.Error().Err(fmt.Errorf("PANIC: %v", err)).Send()
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (s *Server) Start() error {
	if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.http.Shutdown(context.Background())
}
