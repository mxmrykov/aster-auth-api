package main

import (
	"github.com/mxmrykov/asterix-auth/internal/config"
	"github.com/mxmrykov/asterix-auth/internal/service"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Timestamp().Msg("starting auth api...")

	log.Info().Timestamp().Msg("initializing config...")

	cfg, logger, err := config.InitConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize config")
	}

	log.Info().Timestamp().Msg("initializing service...")

	s, err := service.NewService(cfg, logger)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize service")
	}

	log.Info().Timestamp().Msg("starting service...")

	go func() {
		if err = s.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start service")
		}
	}()

	<-utils.GracefulShutDown()

	log.Info().Timestamp().Msg("graceful shutdown")

	if err = s.Stop(); err != nil {
		log.Fatal().Err(err).Msg("failed to stop service")
	}

	log.Info().Timestamp().Msg("service stopped")
}
