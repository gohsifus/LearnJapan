package models

import "time"

// Session Сессии пользователей
type Session struct {
	Id      string    `gorm:"column:id"`
	UserId  int       `gorm:"column:user_id"`
	Expires time.Time `gorm:"column:expires"`
}

func (s Session) TableName() string {
	return "production.sessions"
}
