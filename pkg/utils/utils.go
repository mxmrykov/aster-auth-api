package utils

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mxmrykov/asterix-auth/internal/model"
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

func ValidateAsidToken(token, signature string) (model.SidToken, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&model.SidToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signature), nil
		},
	)

	if claims, ok := parsedToken.Claims.(*model.SidToken); ok && parsedToken.Valid {
		return *claims, nil
	}

	return model.SidToken{}, err
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

func ValidateXTempAuthToken(XAuthToken, signature string) (model.XAuthToken, error) {
	parsedToken, err := jwt.ParseWithClaims(
		XAuthToken,
		&model.XAuthToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signature), nil
		},
	)

	if claims, ok := parsedToken.Claims.(*model.XAuthToken); ok && parsedToken.Valid {
		return *claims, nil
	}

	return model.XAuthToken{}, err
}
