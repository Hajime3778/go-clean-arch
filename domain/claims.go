package domain

import "github.com/form3tech-oss/jwt-go"

type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"name"`
	jwt.StandardClaims
}
