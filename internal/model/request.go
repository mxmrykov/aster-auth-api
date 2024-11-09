package model

type ExtAuthRequest struct {
	Login string `json:"login" binding:"required"`
}
