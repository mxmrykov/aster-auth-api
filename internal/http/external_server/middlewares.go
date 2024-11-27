package external_server

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/pkg/responize"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
)

func (s *Server) internalAuthMiddleWare(ctx *gin.Context) {
	authToken := ctx.GetHeader("X-Auth-Token")

	if authToken == "" {
		s.svc.Log().Error().Msg("Empty auth token")
		responize.R(ctx, nil, http.StatusBadRequest, "Empty auth token", true)
		ctx.Abort()
		return
	}

	jwtSecret, err := s.svc.VaultGetter().GetSecret(
		ctx,
		s.svc.CfgGetter().Vault.TokenRepo.Path,
		s.svc.CfgGetter().Vault.TokenRepo.AppJwtSecretName,
	)

	if err != nil {
		s.svc.Log().Error().Err(err).Msg("vault error")
		responize.R(ctx, nil, http.StatusInternalServerError, "Internal authorization error", true)
		ctx.Abort()
		return
	}

	authPayload, err := utils.ValidateAsidToken(authToken, jwtSecret)

	if err != nil {
		s.svc.Log().Error().Err(err).Msg("token error")
		responize.R(ctx, nil, http.StatusBadRequest, "Invalid token", true)
		ctx.Abort()
		return
	}

	ctx.Set("iaid", authPayload.Iaid)
	ctx.Set("asid", authPayload.Asid)
	ctx.Next()
}

func (s *Server) footPrintAuth(ctx *gin.Context) {
	c, err := ctx.Request.Cookie("X-Client-Footprint")

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			ctx.Set("X-Client-Footprint", s.setFootprintCookie(ctx))
		default:
			responize.R(ctx, nil, http.StatusBadRequest, "Footprint: signature failed", true)
			ctx.Abort()
			return
		}
	} else {
		fpClient := s.svc.CacheGetter().GetClient(c.Value)

		switch {
		case c.Value == "":
			ctx.Set("X-Client-Footprint", s.setFootprintCookie(ctx))
			ctx.Next()
		case fpClient == nil:
			s.dropFootprintCookie(ctx)

			responize.R(ctx, nil, http.StatusUnauthorized, "Footprint: invalid signature", true)
			ctx.Abort()
			return
		case fpClient.RateLimitRemain <= 1:
			if !time.Now().After(fpClient.LastReq.Add(s.svc.CfgGetter().ExternalServer.RateLimiterTimeframe)) {
				fpClient.LastReq = time.Now()
				fpClient.LastUpdated = time.Now()
				s.svc.CacheGetter().SetClient(c.Value, fpClient)

				responize.R(ctx, nil, http.StatusTooManyRequests, "Footprint: rate limited", true)
				ctx.Abort()
				return
			}

			fpClient.RateLimitRemain = s.svc.CfgGetter().ExternalServer.RateLimiterCap
		default:
			fpClient.RateLimitRemain -= 1
		}

		fpClient.LastReq = time.Now()
		fpClient.LastUpdated = time.Now()

		s.svc.CacheGetter().SetClient(c.Value, fpClient)
		ctx.Set("X-Client-Footprint", c.Value)
	}

	ctx.Next()
}

func (s *Server) dropFootprintCookie(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "X-Client-Footprint",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) setFootprintCookie(ctx *gin.Context) string {
	sign := base64.StdEncoding.EncodeToString(
		[]byte(
			uuid.New().String(),
		),
	)

	s.svc.CacheGetter().SetClient(strings.ToUpper(sign), &cache.Props{
		RateLimitRemain: 5,
		LastReq:         time.Now(),
		LastUpdated:     time.Now(),
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "X-Client-Footprint",
		Value:    strings.ToUpper(sign),
		Path:     "/",
		MaxAge:   s.svc.CfgGetter().ExternalServer.RateLimitCookieLifetime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return strings.ToUpper(sign)
}
