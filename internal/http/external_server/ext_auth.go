package external_server

import (
	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/pkg/utils"
	"net/http"
)

func (s *Server) authorizeExternal(ctx *gin.Context) {
	utils.Responize(
		ctx,
		s.cache.MapAllCl(),
		http.StatusOK,
		"",
		false,
	)
}
