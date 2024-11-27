package model

type ExtAuthRequest struct {
	Login string `json:"login" binding:"required"`
}

type IntAuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
