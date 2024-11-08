package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mxmrykov/asterix-auth/pkg/logger"
	"github.com/rs/zerolog"
)

type (
	Auth struct {
		UseStackTrace bool `yaml:"useStackTrace"`

		ExternalServer ExternalServer `yaml:"externalServer"`
		GrpcAST        GrpcAST        `yaml:"grpcAST"`
		GrpcOAuth      GrpcOAuth      `yaml:"grpcOAuth"`
		Vault          Vault          `yaml:"vault"`
	}

	ExternalServer struct {
		Port                 string        `yaml:"port"`
		RateLimiterTimeframe time.Duration `yaml:"rateLimiterTimeframe"`
		RateLimiterCap       uint8         `yaml:"rateLimiterCap"`
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

			AppJwtSecretName   string `yaml:"appJwtSecret"`
			AstJwtSecretName   string `yaml:"astJwtSecret"`
			OAuthJwtSecretName string `yaml:"oAuthJwtSecret"`
		} `yaml:"tokenRepo"`
	}
)

func InitConfig() (*Auth, *zerolog.Logger, error) {
	cfg := *new(Auth)

	if os.Getenv("BUILD_ENV") == "" {
		return nil, nil, errors.New("build environment is not assigned")
	}

	path := fmt.Sprintf("./deploy/%s.yaml", os.Getenv("BUILD_ENV"))

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, nil, err
	}

	l := logger.NewLogger(cfg.UseStackTrace)

	return &cfg, l, nil
}
