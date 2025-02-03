package token

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	signingKey string
}

func New(signKey string) *Token {
	return &Token{
		signingKey: signKey,
	}
}

var (
	//errParseToken   = errors.New("token parse message")
	errInvalidToken = errors.New("invalid token")
	errTokenClams   = errors.New("couldn't parse claims")
	errTokenExpired = errors.New("token is expired")
)

func KeysExtract(jwtToken string, tokenSecret string) (int, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil //errParseToken
		}
		return []byte(tokenSecret), nil
	})

	if token == nil || err != nil {
		return 0, errInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errTokenClams
	}

	exp := claims["exp"].(float64)                    //дата, до которой действует ключ
	userID, _ := strconv.Atoi(claims["sub"].(string)) //id пользователя, под которым логинились
	if int64(exp) < time.Now().Local().Unix() {
		return 0, errTokenExpired
	}

	return userID, err
}

func (t Token) NewJWT(userId int, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   strconv.Itoa(userId),
	})

	tokenString, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		println(err)
	}

	return tokenString, nil
}

func (t Token) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
