package storage_test

import (
	"medods-internship/config"
	"medods-internship/utils"
	"testing"
	"medods-internship/models"  
	"medods-internship/storage" 
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, error) {
	appConfig, err := config.LoadEnv()
	utils.Log(err)
	db, err := config.ConnectToDB(appConfig)
	db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	return db, err
}

// тестируем сохранение токена в базе данных
func TestSaveRefreshToken(t *testing.T) {
	// коннектимся к датабазе
	db, err := setupDB()
	assert.NoError(t, err)
	userStorage := storage.NewStorage(db)

	// создаем данные, которые будем хранить в базе
	userID := "1098754362"
	hashedRefreshToken := "ajksdfgcuqygviruwyegfaksjdfasfc"

	// сохраняем и проверяем что нет ошибки
	err = userStorage.SaveRefreshToken(userID, hashedRefreshToken)
	assert.NoError(t, err)

	// проверяем в самой базе данных что токен сохранился
	var refreshToken models.RefreshToken
	result := db.Where("user_id = ?", userID).First(&refreshToken)

	assert.NoError(t, result.Error)
	assert.Equal(t, userID, refreshToken.UserID)
	assert.Equal(t, hashedRefreshToken, refreshToken.Token)

}

// тестируем получение токена из базы
func TestGetRefreshToken(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	userStorage := storage.NewStorage(db)

	// берем данные из предыдущего теста 
	userID := "1098754362"
	hashedRefreshToken := "ajksdfgcuqygviruwyegfaksjdfasfc"

	token, err := userStorage.GetRefreshToken(userID)

	assert.NoError(t, err)
	assert.Equal(t, hashedRefreshToken, token)
}