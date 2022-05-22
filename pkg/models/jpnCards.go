package models

import(
	"fmt"
)

//GetCardById Возвращает карточку по id
func GetCardById(id string) JpnCards{
	rows, err := DB.Query("select * from Cards where Id = ?", id)
	if err != nil{
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	card := JpnCards{}
	if(rows.Next()){
		rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
	}

	fmt.Println(id)
	fmt.Println(card)

	return card
}

//GetCardList Возвращает все карточки
func GetCardList() []JpnCards{
	rows, err := DB.Query("select * from Cards")
	if err != nil{
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	cards := []JpnCards{}

	for rows.Next(){
		card := JpnCards{}
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		cards = append(cards, card)
	}

	return cards
}

//Add добавляет карточку
func (card *JpnCards) Add() (bool, error){
	sql := "INSERT INTO cards(inJapan, inRussian, mark, dateAdd) VALUES(?, ?, ?, ?)"
	_, err := DB.Exec(sql, card.InJapan, card.InRussian, card.Mark, card.DateAdd)
	if err != nil{
		return false, err
	}

	return true, nil
}
