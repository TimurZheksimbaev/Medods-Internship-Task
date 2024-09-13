package services

import (

	"crypto/rand"
	"encoding/base64"
	"medods-internship/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	appConfig *config.AppConfig
}

func NewTokenService(appConfig *config.AppConfig) *TokenService {
	return &TokenService{
		appConfig: appConfig,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.StandardClaims
}

func (ts *TokenService) GenerateAccessToken(userID, ip string) (string, error) {
	expirationTime := time.Now().Add(ts.appConfig.TokenExpiration * time.Minute)
	claims := &Claims{
		UserID: userID,
		IP:     ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(ts.appConfig.TokenSecret))
}

func (ts *TokenService) GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}


