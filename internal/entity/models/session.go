package models

// Session Сессии пользователей
type Session struct {
	SessionId string
	UserId    int
	Expires   string
}
