package external_server

import (
	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/internal/model"
	"github.com/mxmrykov/asterix-auth/pkg/responize"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
	"net/http"
)

func (s *Server) authorizeExternal(ctx *gin.Context) {
	r := new(model.ExtAuthRequest)

	if e := ctx.ShouldBindJSON(r); e != nil {
		responize.R(ctx, nil, http.StatusBadRequest, "No login provided", true)
		ctx.Abort()
	}

	cc, err := s.svc.VaultGetter().GetSecret(s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AstJwtSecretName)

	if err != nil {
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
	}

	has, iaid, asid, err := s.svc.GrpcAstGetter().GetIAID(ctx, r.Login, cc)

	switch {
	case err != nil:
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
	case !has:
		responize.R(ctx, nil, http.StatusBadRequest, "No such user", true)
		ctx.Abort()
	}

	appSecret, err := s.svc.VaultGetter().GetSecret(s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AppJwtSecretName)

	if err != nil {
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
	}

	assignedToken, err := utils.AssignSidToken(iaid, asid, appSecret)

	responize.R(
		ctx,
		model.ExtAuthResponse{
			SidToken: assignedToken,
		},
		http.StatusOK,
		"",
		false,
	)
}
