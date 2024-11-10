package oauth

import "github.com/mxmrykov/asterix-auth/internal/config"

type IOAuth interface {
}

type OAuth struct {
}

func NewGrpcOAuthClient(cfg *config.GrpcOAuth) (IOAuth, error) {
	return &OAuth{}, nil
}
