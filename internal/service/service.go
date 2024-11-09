package service

import (
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/mxmrykov/asterix-auth/internal/grpc-client/ast"
	"github.com/mxmrykov/asterix-auth/internal/grpc-client/oauth"
	"github.com/mxmrykov/asterix-auth/internal/http/external_server"
	"github.com/mxmrykov/asterix-auth/pkg/clients/vault"
	"github.com/rs/zerolog"
)

type IService interface {
	Start() error
	Stop() error
}

type Service struct {
	Zerolog *zerolog.Logger
	Cfg     *config.Auth

	Cache cache.ICache

	Server *external_server.Server

	Vault vault.IVault

	GrpcAst   ast.IAst
	GrpcOAuth oauth.IOAuth
}

func NewService(cfg *config.Auth, logger *zerolog.Logger) (IService, error) {
	v, err := vault.NewVault(&cfg.Vault)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing vault client")
	}

	grpcAst, err := ast.NewAst(&cfg.GrpcAST)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing AST GRPC client")
	}

	grpcOAuth, err := oauth.NewGrpcOAuthClient(&cfg.GrpcOAuth)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing OAuth GRPC client")
	}

	c := cache.NewCache()

	svc := &Service{
		Zerolog:   logger,
		Cfg:       cfg,
		Cache:     c,
		Vault:     v,
		GrpcAst:   grpcAst,
		GrpcOAuth: grpcOAuth,
	}

	svc.Server = external_server.NewServer(logger, svc)

	return svc, nil
}

func (s *Service) Start() error {
	return s.Server.Start()
}
func (s *Service) Stop() error {
	return s.Server.Stop()
}
func (s *Service) VaultGetter() vault.IVault {
	return s.Vault
}
func (s *Service) GrpcAstGetter() ast.IAst {
	return s.GrpcAst
}
func (s *Service) GrpcOAuthGetter() oauth.IOAuth {
	return s.GrpcOAuth
}
func (s *Service) CacheGetter() cache.ICache { return s.Cache }
func (s *Service) CfgGetter() *config.Auth   { return s.Cfg }
func (s *Service) Log() *zerolog.Logger      { return s.Zerolog }
