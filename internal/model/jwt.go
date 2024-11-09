package model

import "github.com/golang-jwt/jwt"

type SidToken struct {
	Iaid          string `json:"IAID"`
	Asid          string `json:"ASID"`
	SignatureDate string `json:"signatureDate"`
	jwt.StandardClaims
}
