package configs

import (
	"LearnJapan.com/constants"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	Yandex
	MYSQL
}

type Yandex struct {
	YandexApiUrl      string
	YandexApiToken    string
	YandexApiFolderId string
}

type MYSQL struct {
	DB       string `envconfig:"MYSQL_DATABASE"`
	User     string `envconfig:"MYSQL_USER"`
	Password string `envconfig:"MYSQL_PASSWORD"`
	Host     string `envconfig:"MYSQL_HOST"`
	Port     string `envconfig:"MYSQL_PORT"`
}

func NewConfigs() (*Configs, error) {
	if err := godotenv.Load(constants.PATH_ENV_FILE); err != nil {
		return nil, err
	}

	config := &Configs{}

	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	return config, nil
}
