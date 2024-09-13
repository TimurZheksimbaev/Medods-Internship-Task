package controllers

import (
	"fmt"
	"medods-internship/config"
	"medods-internship/services"
	"medods-internship/storage"
	"medods-internship/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
	storage *storage.UserStorage
	tokenService *services.TokenService
	appConfig *config.AppConfig
}

func NewTokenController(storage *storage.UserStorage, tokenService *services.TokenService, appConfig *config.AppConfig) *TokenController {
	return &TokenController{
		storage: storage,
		tokenService: tokenService,
		appConfig: appConfig,
	}
}

func (tc *TokenController) GenerateTokens(c *gin.Context) {
	// достаем user_id и ip из запроса
	userID := c.Query("user_id")
	ip := c.ClientIP()

	// генерируем access и refresh токены
	accessToken, err := tc.tokenService.GenerateAccessToken(userID, ip)
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := tc.tokenService.GenerateRefreshToken()
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	// хэшируем refresh токен
	hashedRefreshToken, err := utils.HashRefreshToken(refreshToken)
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash refresh token"})
		return
	}

	// сохраняем refresh токен
	err = tc.storage.SaveRefreshToken(userID, hashedRefreshToken)
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
		return
	}

	// возвращаем access и refresh токены
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (tc *TokenController) RefreshTokens(c *gin.Context) {
	var request struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	// достаем access и refresh токены из запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Log(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// проверяем валидность access токена
	claims := &services.Claims{}
	accessToken, err := jwt.ParseWithClaims(request.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tc.appConfig.TokenSecret), nil
	})
	utils.Log(err)

	if err != nil || !accessToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	// проверяем валидность refresh токена
	hashedRefreshToken, err := tc.storage.GetRefreshToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get refresh token"})
		return
	}

	if !utils.CompareRefreshTokens(hashedRefreshToken, request.RefreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// IP пользователся изменился, поэтому отправляем ему сообщение на почту
	if claims.IP != c.ClientIP() {
		// Send email warning
		utils.SendEmailWarning(claims.UserID)
	}

	// генерируем новые access и refresh токены
	newAccessToken, err := tc.tokenService.GenerateAccessToken(claims.UserID, c.ClientIP())
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return
	}

	newRefreshToken, err := tc.tokenService.GenerateRefreshToken()
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new refresh token"})
		return
	}

	newHashedRefreshToken, err := utils.HashRefreshToken(newRefreshToken)
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash new refresh token"})
		return
	}

	err = tc.storage.SaveRefreshToken(claims.UserID, newHashedRefreshToken)
	utils.Log(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save new refresh token"})
		return
	}

	// возвращаем новые access и refresh токены
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}