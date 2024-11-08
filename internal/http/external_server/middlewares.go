package external_server

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mxmrykov/asterix-auth/internal/cache"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
	"net/http"
	"strings"
	"time"
)

func (s *Server) internalAuthMiddleWare(ctx *gin.Context) {

}

func (s *Server) footPrintAuth(ctx *gin.Context) {
	c, err := ctx.Request.Cookie("X-Client-Footprint")

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			s.setFootprintCookie(ctx)
		default:
			utils.Responize(ctx, nil, http.StatusBadRequest, "Footprint: signature failed", true)
			ctx.Abort()
			return
		}
	} else {
		fpClient := s.cache.GetClient(c.Value)

		switch {
		case c.Value == "":
			s.setFootprintCookie(ctx)
			ctx.Next()
		case fpClient == nil:
			utils.Responize(ctx, nil, http.StatusUnauthorized, "Footprint: invalid signature", true)
			ctx.Abort()
			return
		case fpClient.RateLimitRemain <= 0:
			if !time.Now().After(fpClient.LastReq.Add(5 * time.Second)) {
				utils.Responize(ctx, nil, http.StatusTooManyRequests, "Footprint: rate limited", true)
				ctx.Abort()
				return
			}

			fpClient.RateLimitRemain = 5
		default:
			fpClient.RateLimitRemain -= 1
		}

		fpClient.LastUpdated = time.Now()
		fpClient.LastReq = time.Now()

		s.cache.SetClient(c.Value, fpClient)
	}

	ctx.Next()
}

func (s *Server) setFootprintCookie(ctx *gin.Context) {
	sign := base64.StdEncoding.EncodeToString(
		[]byte(
			uuid.New().String(),
		),
	)

	cookie := http.Cookie{
		Name:     "X-Client-Footprint",
		Value:    strings.ToUpper(sign),
		Path:     "/",
		MaxAge:   900,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	s.cache.SetClient(strings.ToUpper(sign), &cache.Props{
		RateLimitRemain: 5,
		LastReq:         time.Now(),
		LastUpdated:     time.Now(),
	})

	http.SetCookie(ctx.Writer, &cookie)
}
