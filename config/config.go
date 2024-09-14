package config

import (
	"errors"

	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	DatabaseHost string `mapstructure:"DATABASE_HOST"`
	DatabaseUsername string `mapstructure:"DATABASE_USERNAME"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`
	DatabasePort int `mapstructure:"DATABASE_PORT"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	TokenExpiration time.Duration `mapstructure:"TOKEN_EXPIRATION"`
	TokenSecret string `mapstructure:"TOKEN_SECRET"`
}


func LoadEnv() (*AppConfig, error) {
	projectRoot, err := filepath.Abs(".")
	envPath := filepath.Join(projectRoot, "app.env")
	viper.SetConfigFile(envPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.New("Could not read config file")
	}
	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	return &appConfig, err
}

func ConnectToDB(config *AppConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
}