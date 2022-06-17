package configs

import (
	"log"
	"os"
)

var Cfg Configs

func init(){
	url, ok := os.LookupEnv("YANDEX_API_URL")
	token, ok := os.LookupEnv("YANDEX_API_TOKEN")
	folderId, ok := os.LookupEnv("YANDEX_API_FOLDER_ID")
	dbConnString, ok := os.LookupEnv("DB_CONNECTION_STRING")

	if !ok{
		log.Fatalln("Error: Переменные окружения не инициализированны")
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
