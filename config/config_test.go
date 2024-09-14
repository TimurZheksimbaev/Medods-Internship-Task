package config_test

import (
	"medods-internship/config"
	"os"
	"path/filepath"

	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var projectRoot, _ = filepath.Abs("../") 
var err = os.Chdir(projectRoot)

// проверяем что ссылка на датабазу непустая и проверяем подключение
func TestConnectToDB(t *testing.T) {
	appConfig, _ := config.LoadEnv()
	assert.NotEmpty(t, appConfig.DatabaseUrl)

	db, err := config.ConnectToDB(appConfig)
	assert.NoError(t, err)
	assert.IsType(t, &gorm.DB{}, db)
}


