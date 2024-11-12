package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/mxmrykov/asterix-auth/internal/model"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutDown() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return c
}

func AssignAsidToken(iaid, asid, signature string) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.SidToken{
		Iaid:          iaid,
		Asid:          asid,
		SignatureDate: time.Now().Format(time.RFC3339),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "aster-auth",
		},
	})
	return unsignedToken.SignedString([]byte(signature))
}

func AssignXAuthToken(asid, signature string) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.XAuthToken{
		Asid:          asid,
		SignatureDate: time.Now().Format(time.RFC3339),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "aster-auth",
		},
	})
	return unsignedToken.SignedString(signature)
}
