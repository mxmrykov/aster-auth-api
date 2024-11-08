package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/asterix-auth/internal/model"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutDown() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return c
}

func Responize(ctx *gin.Context, data interface{}, code int, message string, error bool) {
	ctx.JSON(code, model.Response{
		Payload: data,
		Status:  code,
		Message: message,
		Error:   error,
	})
}
