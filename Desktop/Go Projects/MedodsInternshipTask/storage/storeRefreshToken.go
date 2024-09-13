package storage

import (
	"medods-internship/models"
	"medods-internship/utils"

	"gorm.io/gorm"
)


type UserStorage struct {
	DB *gorm.DB
}

func NewStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{DB: db}
}

// сохраняем рефреш токен
func (us *UserStorage) SaveRefreshToken(userID, hashedRefreshToken string) error {
	refreshToken := models.RefreshToken{
		UserID: userID,
		Token:  hashedRefreshToken,
	}
	result := us.DB.Create(&refreshToken)
	if result.Error != nil {
		return result.Error
	}
	utils.LogMessage("Successfully saved refresh token")
	return nil
}

// получаем рефреш токен
func (us *UserStorage) GetRefreshToken(userID string) (string, error) {
	var token models.RefreshToken
	result := us.DB.Where("user_id = ?", userID).First(&token)
	if result.Error != nil {
		return "", result.Error
	}
	return token.Token, nil
}