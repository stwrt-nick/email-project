package model

type AddPostRequest struct {
}

type AddPostResponse struct {
	Status string `json:"status"`
	Err    string `json:"err"`
}

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}
