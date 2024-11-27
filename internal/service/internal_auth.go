package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/internal/grpc-client/oauth"
	"github.com/mxmrykov/asterix-auth/internal/model"
	"github.com/mxmrykov/asterix-auth/pkg/responize"
)

func (s *Service) InternalAuthorize(ctx *gin.Context, login, pass string) (*model.ClientResponse, error) {
	cc, err := s.VaultGetter().GetSecret(ctx, s.CfgGetter().Vault.TokenRepo.Path, s.CfgGetter().Vault.TokenRepo.OAuthJwtSecretName)

	if err != nil {
		s.Log().Err(err).Msg("Failed to get secret")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return nil, err
	}

	clientData, err := s.GrpcOAuth.Authorize(ctx, oauth.Req{
		Login:       login,
		Password:    pass,
		IAID:        ctx.GetString("iaid"),
		ConfirmCode: cc,
		ASID:        ctx.GetString("asid"),
	})

	if err != nil {
		s.Log().Err(err).Msg("Failed to authorize client")
		return nil, err
	}

	return &model.ClientResponse{
		ClientID:     clientData.ClientID,
		ClientSecret: clientData.ClientSecret,
		OAuthCode:    clientData.OAuthCode,
	}, nil
}
