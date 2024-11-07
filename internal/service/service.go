package service

import (
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/mxmrykov/asterix-auth/internal/grpc"
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

	Server external_server.IServer

	Vault vault.IVault

	GrpcAst   grpc.IAst
	GrpcOAuth grpc.IOAuth
}

func NewService(cfg *config.Auth, logger *zerolog.Logger) (IService, error) {
	v, err := vault.NewVault(&cfg.Vault)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing vault client")
	}

	grpcAst, err := grpc.NewGrpcAstClient(&cfg.GrpcAST)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing AST GRPC client")
	}

	grpcOAuth, err := grpc.NewGrpcOAuthClient(&cfg.GrpcOAuth)

	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing OAuth GRPC client")
	}

	c := cache.NewCache()

	return &Service{
		Zerolog:   logger,
		Cfg:       cfg,
		Cache:     c,
		Vault:     v,
		Server:    external_server.NewServer(&cfg.ExternalServer, logger, c),
		GrpcAst:   grpcAst,
		GrpcOAuth: grpcOAuth,
	}, nil
}

func (s *Service) Start() error {
	return s.Server.Start()
}
func (s *Service) Stop() error {
	return s.Server.Stop()
}
