package models

// User Пользователь
type User struct {
	Id       int    `gorm:"column:id"`
	Email    string `gorm:"column:email"`
	Login    string `gorm:"column:login"`
	Password string `gorm:"column:password"`
}

func (u User) TableName() string {
	return "production.users"
}
