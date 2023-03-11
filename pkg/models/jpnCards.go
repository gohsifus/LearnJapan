package models

import (
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

//GetCardListBySessionId Вернет карточки принадлежащие конкретному пользователю по id действуещей сессии
func GetCardListBySessionId(sessionId string) ([]JpnCards, error){
	sql := "SELECT cards.* " +
			"FROM cards " +
			"LEFT JOIN sessions " +
			"	ON cards.userId = sessions.userId " +
			"WHERE sessionId = ? and expires > NOW();"

	rows, err := DB.Query(sql, sessionId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	cards := []JpnCards{}

	for rows.Next(){
		card := JpnCards{}
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		cards = append(cards, card)
	}

	return cards, nil
}

//Add добавляет карточку
func (card *JpnCards) Add() (bool, error){
	sql := `INSERT INTO cards(inJapan, inRussian, mark, dateAdd, userId) select ?, ?, ?, ?, ? 
			WHERE NOT EXISTS (SELECT 1 FROM cards WHERE inJapan = ?)`
	_, err := DB.Exec(sql, card.InJapan, card.InRussian, card.Mark, card.DateAdd, card.UserId, card.InJapan)
	if err != nil{
		return false, err
	}

	return true, nil
}

//GetRandCard Вернет случайную карточку
func GetRandCard() (JpnCards, error){
	sql := "SELECT * FROM cards ORDER BY rand() LIMIT 1"

	rows, err := DB.Query(sql)
	if err != nil{
		return JpnCards{}, err
	}
	defer rows.Close()

	randCard := JpnCards{}
	if rows.Next(){
		rows.Scan(
			&randCard.Id,
			&randCard.InJapan,
			&randCard.InRussian,
			&randCard.Mark,
			&randCard.DateAdd,
			&randCard.UserId)
	}

	return randCard, nil
}

//GetRandCardForUser Вернет случайную карточку пользователя
func GetRandCardForUser(sessionId string) (JpnCards, error) {
	sql := "SELECT cards.* " +
		   "FROM cards " +
		   "LEFT JOIN sessions " +
		   "	ON cards.userId = sessions.userId " +
		   "WHERE sessions.sessionId = ? AND expires > NOW()" +
		   "ORDER BY rand() LIMIT 1"


	rows, err := DB.Query(sql, sessionId)
	defer rows.Close()
	if err != nil{
		return JpnCards{}, err
	}

	card := JpnCards{}

	if rows.Next() {
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		if err != nil{
			return JpnCards{}, err
		}
	}

	return card, nil
}

//GetCardByInJapan Вернет карточку по слову на японском
func GetCardByInJapan(inJapan string) (JpnCards, bool){
	sql := "SELECT * FROM cards WHERE InJapan = ?"
	rows, err := DB.Query(sql, inJapan)
	defer rows.Close()
	if err != nil {
		return JpnCards{}, false
	}

	resultCard := JpnCards{}
	if rows.Next(){
		rows.Scan(
			&resultCard.Id,
			&resultCard.InJapan,
			&resultCard.InRussian,
			&resultCard.Mark,
			&resultCard.DateAdd,
			&resultCard.UserId)
	} else {
		return JpnCards{}, false
	}

	return resultCard, true
}

//UpdateCardMark Изменит значение оценки слова на mark
func UpdateCardMark(id, mark int) error{
	sql := "UPDATE cards SET mark = mark + ? WHERE id = ?"
	_, err := DB.Exec(sql, mark, id)
	if err != nil{
		return err
	}
	return nil
}
