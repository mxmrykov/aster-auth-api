package external_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/internal/model"
	"github.com/mxmrykov/asterix-auth/pkg/responize"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
)

func (s *Server) authorizeExternal(ctx *gin.Context) {
	r := new(model.ExtAuthRequest)

	if err := ctx.ShouldBindJSON(r); err != nil {
		s.svc.Log().Err(err).Msg("Error in binding JSON")
		responize.R(ctx, nil, http.StatusBadRequest, "No login provided", true)
		ctx.Abort()
		return
	}

	cc, err := s.svc.VaultGetter().GetSecret(ctx, s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AstJwtSecretName)

	if err != nil {
		s.svc.Log().Err(err).Msg("Failed to get secret")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return
	}

	has, iaid, asid, err := s.svc.GrpcAstGetter().GetIAID(ctx, r.Login, cc)

	switch {
	case err != nil:
		s.svc.Log().Err(err).Msg("gRPC request failed")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return
	case !has:
		s.svc.Log().Err(err).Msg("No such user")
		responize.R(ctx, nil, http.StatusBadRequest, "No such user", true)
		ctx.Abort()
		return
	}

	appSecret, err := s.svc.VaultGetter().GetSecret(ctx, s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AppJwtSecretName)

	if err != nil {
		s.svc.Log().Err(err).Msg("Failed to get app secret")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return
	}

	assignedToken, err := utils.AssignAsidToken(iaid, asid, appSecret)

	responize.R(
		ctx,
		model.ExtAuthResponse{
			SidToken: assignedToken,
		},
		http.StatusOK,
		"",
		false,
	)

	s.svc.Log().Info().Msg("Successfully external authorized user")
}

func (s *Server) authorizeInternal(ctx *gin.Context) {

}

func (s *Server) checkLogin(ctx *gin.Context) {
	r := new(model.ExtAuthRequest)

	if err := ctx.ShouldBindJSON(r); err != nil {
		s.svc.Log().Err(err).Msg("Error in binding JSON")
		responize.R(ctx, nil, http.StatusBadRequest, "No login provided", true)
		ctx.Abort()
		return
	}

	cc, err := s.svc.VaultGetter().GetSecret(ctx, s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AstJwtSecretName)

	if err != nil {
		s.svc.Log().Err(err).Msg("Failed to get secret")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return
	}

	appSecret, err := s.svc.VaultGetter().GetSecret(ctx, s.svc.CfgGetter().Vault.TokenRepo.Path, s.svc.CfgGetter().Vault.TokenRepo.AppJwtSecretName)

	if err != nil {
		s.svc.Log().Err(err).Msg("Failed to get app secret")
		responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
		ctx.Abort()
		return
	}

	var assignedToken string
	cookie, ok := ctx.Get("X-Client-Footprint")

	if !ok {
		s.svc.Log().Err(err).Msg("X-Client-Footprint is not assigned")
		responize.R(ctx, nil, http.StatusInternalServerError, "X-Client-Footprint is not assigned", true)
		ctx.Abort()
		return
	}

	if v := s.svc.CacheGetter().GetClient(cookie.(string)); v == nil || v.ASID == "" {
		has, _, asid, err := s.svc.GrpcAstGetter().GetIAID(ctx, r.Login, cc)

		switch {
		case err != nil:
			s.svc.Log().Err(err).Msg("gRPC request failed")
			responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
			ctx.Abort()
			return
		case has:
			s.svc.Log().Err(err).Msg("Login is already in use")
			responize.R(ctx, model.CheckLoginResponse{Unused: false}, http.StatusBadRequest, "Login is already in use", true)
			ctx.Abort()
			return
		}

		assignedToken, err = utils.AssignAsidToken("", asid, appSecret)

		if err != nil {
			s.svc.Log().Err(err).Msg("Failed to get assigned token")
			responize.R(ctx, nil, http.StatusInternalServerError, "Confirm login error", true)
			ctx.Abort()
			return
		}

		s.svc.CacheGetter().SetClient(cookie.(string), &cache.Props{
			ASID:            assignedToken,
			RateLimitRemain: v.RateLimitRemain,
			LastReq:         v.LastReq,
			LastUpdated:     v.LastUpdated,
		})
	} else if _, err = utils.ValidateXTempAuthToken(v.ASID, appSecret); err == nil {
		assignedToken = v.ASID
	}

	responize.R(
		ctx,
		model.CheckLoginResponse{
			Unused:         true,
			XTempauthToken: assignedToken,
		},
		http.StatusOK,
		"",
		false,
	)

	s.svc.Log().Info().Msg("Successfully check if login used")
}
