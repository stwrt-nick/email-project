package model

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}
