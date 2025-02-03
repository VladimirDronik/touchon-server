package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type Tokens struct {
	AccessToken  string `json:"api_access_token"`
	RefreshToken string `json:"api_refresh_token"`
}
