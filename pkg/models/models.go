package models

import (
	"database/sql"
)

var DB *sql.DB

//JpnCards Карточка с словом на разных языках
type JpnCards struct{
	Id        int
	InJapan   string
	InRussian string
	Mark      int
	DateAdd   string
	UserId 	  int
}

//User Пользователь
type User struct{
	Id       int
	Email    string
	Login    string
	Password string
}

//Session Сессии пользователей
type Session struct{
	SessionId string
	UserId int
	Expires string
}
