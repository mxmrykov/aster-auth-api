package ast

import (
	"context"
	"fmt"
	"github.com/mxmrykov/asterix-auth/internal/config"
	ast "github.com/mxmrykov/asterix-auth/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type IAst interface {
	GetIAID(ctx context.Context, login, cc string) (bool, string, string, error)
}

type Ast struct {
	Conn        ast.AstClient
	MaxPollTime time.Duration
}

func NewAst(cf *config.GrpcAST) (*Ast, error) {
	c, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cf.Host, cf.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}...,
	)

	if err != nil {
		return nil, err
	}

	return &Ast{
		Conn:        ast.NewAstClient(c),
		MaxPollTime: cf.MaxPollTime,
	}, nil
}

func (a *Ast) GetIAID(ctx context.Context, login, cc string) (bool, string, string, error) {
	r := ast.GetIAIDRequest{Login: login, ConfirmCode: cc}

	res, err := a.Conn.GetIAID(ctx, &r)

	if err != nil {
		return false, "", "", err
	}

	if res.Message == "no such login" {
		return false, "", res.ASID, nil
	}

	return res.Has, res.IAID, res.ASID, nil
}
