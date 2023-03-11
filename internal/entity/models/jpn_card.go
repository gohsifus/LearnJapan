package models

// JpnCard Карточка со словом на разных языках
type JpnCard struct {
	Id        int    `gorm:"column:id"`
	InJapan   string `gorm:"column:in_japan"`
	InRussian string `gorm:"column:in_russian"`
	Mark      int    `gorm:"column:mark"`
	DateAdd   string `gorm:"column:date_add"`
	UserId    int    `gorm:"column:user_id"`
}

func (c JpnCard) TableName() string {
	return "production.cards"
}
