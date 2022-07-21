package configs

import (
	"LearnJapan.com/pkg/logger"
	"os"
)

var Cfg Configs

func init(){
	url, ok := os.LookupEnv("YANDEX_API_URL")
	token, ok := os.LookupEnv("YANDEX_API_TOKEN")
	folderId, ok := os.LookupEnv("YANDEX_API_FOLDER_ID")
	dbConnString, ok := os.LookupEnv("DB_CONNECTION_STRING")

	if !ok{
		logger.Print("FatalError: Переменные окружения не инициализированны")
		os.Exit(1)
	}

	Cfg = Configs{
		YandexApiUrl: url,
		YandexApiToken: token,
		YandexApiFolderId: folderId,
		DBConnectionString: dbConnString,
	}
}

type Configs struct{
	YandexApiUrl       string
	YandexApiToken 	   string
	YandexApiFolderId  string
	DBConnectionString string
}
