package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mxmrykov/asterix-auth/pkg/logger"
	"github.com/rs/zerolog"
	"os"
)

type (
	Auth struct {
		UseStackTrace  bool `yaml:"useStackTrace"`
		AppJwtSecret   string
		AstJwtSecret   string
		OAuthJwtSecret string

		ExternalServer ExternalServer `yaml:"externalServer"`
		GrpcAST        GrpcAST        `yaml:"grpcAST"`
		GrpcOAuth      GrpcOAuth      `yaml:"grpcOAuth"`
		Vault          Vault          `yaml:"vault"`
	}

	ExternalServer struct {
		Port string `yaml:"port"`
	}

	GrpcAST struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	GrpcOAuth struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Vault struct {
		AuthToken string

		TokenRepo struct {
			Path string `yaml:"path"`

			AppJwtSecretName   string
			AstJwtSecretName   string
			OAuthJwtSecretName string
		} `yaml:"tokenRepo"`
	}
)

func InitConfig() (*Auth, *zerolog.Logger, error) {
	cfg := *new(Auth)
	path := fmt.Sprintf("./deploy/%s.yaml", os.Getenv("BUILD_ENV"))

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, nil, err
	}

	l := logger.NewLogger(cfg.UseStackTrace)

	return &cfg, l, nil
}
