package models

// JpnCard Карточка с словом на разных языках
type JpnCard struct {
	Id        int
	InJapan   string
	InRussian string
	Mark      int
	DateAdd   string
	UserId    int
}
