package models

func GetById(id int) JpnCards{
	return JpnCards{
		id: id,
		inRussian: "Тест",
		inJapan: "てすと",
		dateAdd: "01.12.2022",
	}
}
