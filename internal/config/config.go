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
		Port                    int           `yaml:"port"`
		RateLimiterTimeframe    time.Duration `yaml:"rateLimiterTimeframe"`
		RateLimiterCap          uint8         `yaml:"rateLimiterCap"`
		RateLimitCookieLifetime int           `yaml:"rateLimitCookieLifetime"`
	}

	GrpcAST struct {
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		MaxPollTime time.Duration `yaml:"maxPollTime"`
	}

	GrpcOAuth struct {
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		MaxPollTime time.Duration `yaml:"maxPollTime"`
	}

	Vault struct {
		AuthToken     string        `env:"VAULT_AUTH_TOKEN"`
		Host          string        `yaml:"host"`
		Port          int           `yaml:"port"`
		ClientTimeout time.Duration `yaml:"clientTimeout"`

		TokenRepo struct {
			Path string `yaml:"path"`

			AppJwtSecretName   string `yaml:"appJwtSecretName"`
			AstJwtSecretName   string `yaml:"astJwtSecretName"`
			OAuthJwtSecretName string `yaml:"oAuthJwtSecretName"`
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

	l.Info().Msgf("%v", cfg)

	return &cfg, l, nil
}
