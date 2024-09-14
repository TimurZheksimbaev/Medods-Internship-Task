package controllers_test

import (
	"bytes"
	"encoding/json"
	"medods-internship/config"
	"medods-internship/controllers"
	"medods-internship/models"
	"medods-internship/services"
	"medods-internship/storage"
	"medods-internship/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var projectRoot, _ = filepath.Abs("../") 
var err = os.Chdir(projectRoot)


func setupTestRouter() *gin.Engine {
	// инициализируем компоненты проекта
	appConfig, err := config.LoadEnv()
	utils.LogExit(err)

	db, err := config.ConnectToDB(appConfig)
	utils.LogExit(err)

	db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	storage := storage.NewStorage(db)
	tokenService := services.NewTokenService(appConfig)

	tokenController := controllers.NewTokenController(storage, tokenService, appConfig)

	r := gin.Default()
	r.POST("/auth/generate-tokens", tokenController.GenerateTokens)
	r.POST("/auth/refresh-tokens", tokenController.RefreshTokens)
	return r
}

// тестируем генерацию токенов
func TestGenerateTokens(t *testing.T) {	
	authRouter := setupTestRouter()

	// формирумем запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/generate-tokens?user_id=test-user", nil)
	req.RemoteAddr = "127.0.0.1:3000"

	authRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// получаем ответ, обрабатываем его
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// тестируем обновление токенов
func TestRefreshToken(t *testing.T) {
   	authRouter := setupTestRouter()

	// сначала генерируем токены
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/generate-tokens?user_id=test-user", nil)
	req.RemoteAddr = "127.0.0.1:3000"

	authRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// получаем токены и парсим их
	var tokenResponse TokenResponse
	json.Unmarshal(w.Body.Bytes(), &tokenResponse)

	assert.NotEmpty(t, tokenResponse.AccessToken)
	assert.NotEmpty(t, tokenResponse.RefreshToken)
	assert.NotEqual(t, tokenResponse.AccessToken, tokenResponse.RefreshToken)

	// ставим невалидные данные
	tokenResponse.RefreshToken = "invalid_token"

	invalidRefreshJSON, _ := json.Marshal(tokenResponse)

	// отправляем запрос на обновление токенов
	req, _ = http.NewRequest("POST", "/auth/refresh-tokens", bytes.NewBuffer(invalidRefreshJSON))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	authRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// если рефреш токен невалидный. то мы не получим новые токены
	var newTokenResponse TokenResponse
	json.Unmarshal(w.Body.Bytes(), &tokenResponse)

	assert.Empty(t, newTokenResponse.AccessToken)
	assert.Empty(t, newTokenResponse.RefreshToken)
}
