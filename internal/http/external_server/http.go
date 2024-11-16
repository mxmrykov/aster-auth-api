package external_server

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/mxmrykov/asterix-auth/internal/grpc-client/ast"
	"github.com/mxmrykov/asterix-auth/internal/grpc-client/oauth"
	"github.com/mxmrykov/asterix-auth/pkg/clients/vault"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type IServer interface {
	Start() error
	Stop() error

	VaultGetter() vault.IVault

	GrpcAstGetter() ast.IAst
	GrpcOAuthGetter() oauth.IOAuth

	CacheGetter() cache.ICache
	CfgGetter() *config.Auth
	Log() *zerolog.Logger
}

type Server struct {
	svc    IServer
	router *gin.Engine
	http   http.Server
}

func NewServer(logger *zerolog.Logger, svc IServer) *Server {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.CustomRecoveryWithWriter(nil, recoveryFunc(logger)),
	)

	s := &Server{
		svc:    svc,
		router: router,
		http: http.Server{
			Addr:    fmt.Sprintf(":%d", svc.CfgGetter().ExternalServer.Port),
			Handler: router,
		},
	}

	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.router.Use(s.footPrintAuth)
	s.router.Use(cors.Default())

	internalAuthGroup := s.router.Group("/auth/api/v1/internal")
	internalAuthGroup.Use(s.internalAuthMiddleWare)
	internalAuthGroup.POST("/new/oauth")

	externalAuthGroup := s.router.Group("/auth/api/v1/external")
	externalAuthGroup.POST("/new/sid", s.authorizeExternal)
	externalAuthGroup.POST("/check/login", s.checkLogin)
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
