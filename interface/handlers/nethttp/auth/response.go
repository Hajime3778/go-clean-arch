package auth

type SignUpResponse struct {
	Token string `json:"token"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
