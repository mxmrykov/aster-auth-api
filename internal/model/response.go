package model

type Response struct {
	Payload interface{} `json:"payload"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
}

type ExtAuthResponse struct {
	SidToken string `json:"sid_token"`
}

type CheckLoginResponse struct {
	Unused         bool   `json:"unused"`
	XTempauthToken string `json:"x_TempAuth_Token"`
}
