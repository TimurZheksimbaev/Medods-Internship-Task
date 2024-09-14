package services_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"medods-internship/config"    
	"medods-internship/services"  
)

var appConfig, _ = config.LoadEnv()
var tokenService = services.NewTokenService(appConfig)

// тестриуем генерацию токена доступа
func TestGenerateAccessToken(t *testing.T) {
	// готовим данные для токена доступа
	userID := "9187546"
	ip := "192.168.0.101"

	// генерируем токен доступа
	tokenString, err := tokenService.GenerateAccessToken(userID, ip)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// парсим токен доступа и проверяем его валидность
	token, err := jwt.ParseWithClaims(tokenString, &services.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(appConfig.TokenSecret), nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, token)

	// проверяем содержимое токена
	if claims, ok := token.Claims.(*services.Claims); ok && token.Valid {
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, ip, claims.IP)
		assert.WithinDuration(t, time.Unix(claims.ExpiresAt, 0), time.Now().Add(appConfig.TokenExpiration*time.Minute), time.Minute)
	} else {
		t.Errorf("Token claims are invalid")
	}
}

// тестируем генерацию рефреш токена
func TestGenerateRefreshToken(t *testing.T) {
	refreshToken, err := tokenService.GenerateRefreshToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	assert.True(t, strings.ContainsAny(refreshToken, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="))
}
