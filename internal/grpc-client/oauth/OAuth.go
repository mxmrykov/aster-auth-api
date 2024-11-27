package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mxmrykov/asterix-auth/internal/config"
	oauth "github.com/mxmrykov/asterix-auth/internal/proto/oauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IOAuth interface {
	Authorize(ctx context.Context, req Req) (*Res, error)
}

type OAuth struct {
	Conn        oauth.OAuthClient
	MaxPollTime time.Duration
}

type Req struct {
	Login       string
	Password    string
	IAID        string
	ConfirmCode string
	ASID        string
}

type Res struct {
	ClientID     string
	ClientSecret string
	OAuthCode    string
	Error        string
}

func NewGrpcOAuthClient(cfg *config.GrpcOAuth) (IOAuth, error) {
	c, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}...,
	)

	if err != nil {
		return nil, err
	}

	return &OAuth{
		Conn:        oauth.NewOAuthClient(c),
		MaxPollTime: cfg.MaxPollTime,
	}, nil
}

func (o *OAuth) Authorize(ctx context.Context, req Req) (*Res, error) {
	r := oauth.AuthorizeRequest{
		Login:       req.Login,
		Password:    req.Password,
		IAID:        req.IAID,
		ConfirmCode: req.ConfirmCode,
		ASID:        req.ASID,
	}

	res, err := o.Conn.Authorize(ctx, &r)

	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	return &Res{
		ClientID:     res.ClientID,
		ClientSecret: res.ClientSecret,
		OAuthCode:    res.OAuthCode,
	}, nil
}
