package grpc

import "github.com/mxmrykov/asterix-auth/internal/config"

type IAst interface {
}

type Ast struct {
}

func NewGrpcAstClient(cfg *config.GrpcAST) (IAst, error) {
	return &Ast{}, nil
}
