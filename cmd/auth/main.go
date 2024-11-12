package main

import (
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/mxmrykov/asterix-auth/internal/service"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, logger, err := config.InitConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize config")
	}

	logger.Info().Timestamp().Msg("config initialized")
	logger.Info().Timestamp().Msg("initializing service...")

	s, err := service.NewService(cfg, logger)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize service")
	}

	logger.Info().Timestamp().Msg("starting service...")

	go func() {
		if err = s.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start service")
		}
	}()

	<-utils.GracefulShutDown()

	logger.Info().Timestamp().Msg("graceful shutdown")

	if err = s.Stop(); err != nil {
		logger.Fatal().Err(err).Msg("failed to stop service")
	}

	logger.Info().Timestamp().Msg("service stopped")
}
