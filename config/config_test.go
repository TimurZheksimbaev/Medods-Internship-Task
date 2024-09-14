package config_test

import (
	"medods-internship/config"

	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)


// проверяем что ссылка на датабазу непустая и проверяем подключение
func TestConnectToDB(t *testing.T) {
	appConfig, _ := config.LoadEnv()
	assert.NotEmpty(t, appConfig.DatabaseUrl)

	db, err := config.ConnectToDB(appConfig)
	assert.NoError(t, err)
	assert.IsType(t, &gorm.DB{}, db)
}

// проверяем что при ложной ссылке подключение к датабазе не получится
func TestConnectToDBFalseUrl(t *testing.T) {
	badConfig := &config.AppConfig{
		DatabaseUrl: "some url",
	}
	db, err := config.ConnectToDB(badConfig)

	assert.Error(t, err)
	assert.Nil(t, db.Error)
}
