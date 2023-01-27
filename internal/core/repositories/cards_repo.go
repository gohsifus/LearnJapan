package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"database/sql"
	"fmt"
)

type CardRepo struct {
	DB *sql.DB
}

func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{
		DB: db,
	}
}

// GetCardById Возвращает карточку по id
func (r CardRepo) GetCardById(id string) models.JpnCard {
	rows, err := r.DB.Query("select * from Cards where Id = ?", id)
	if err != nil {
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	card := models.JpnCard{}
	if rows.Next() {
		rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
	}

	return card
}

// GetCardList Возвращает все карточки
func (r CardRepo) GetCardList() []models.JpnCard {
	rows, err := r.DB.Query("select * from Cards")
	if err != nil {
		panic(fmt.Sprintf("Ошибка выбора даных: %s", err))
	}
	defer rows.Close()

	cards := []models.JpnCard{}

	for rows.Next() {
		card := models.JpnCard{}
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		cards = append(cards, card)
	}

	return cards
}

// GetCardListBySessionId Вернет карточки принадлежащие конкретному пользователю по id действуещей сессии
func (r CardRepo) GetCardListBySessionId(sessionId string) ([]models.JpnCard, error) {
	sql := "SELECT cards.* " +
		"FROM cards " +
		"LEFT JOIN sessions " +
		"	ON cards.userId = sessions.userId " +
		"WHERE sessionId = ? and expires > NOW();"

	rows, err := r.DB.Query(sql, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []models.JpnCard{}

	for rows.Next() {
		card := models.JpnCard{}
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		cards = append(cards, card)
	}

	return cards, nil
}

// Add добавляет карточку
func (r CardRepo) Add(card *models.JpnCard) (bool, error) {
	sql := "INSERT INTO cards(inJapan, inRussian, mark, dateAdd, userId) select ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM cards WHERE inJapan = ?)"

	_, err := r.DB.Exec(sql, card.InJapan, card.InRussian, card.Mark, card.DateAdd, card.UserId, card.InJapan)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRandCard Вернет случайную карточку
func (r CardRepo) GetRandCard() (models.JpnCard, error) {
	sql := "SELECT * FROM cards ORDER BY rand() LIMIT 1"

	rows, err := r.DB.Query(sql)
	if err != nil {
		return models.JpnCard{}, err
	}
	defer rows.Close()

	randCard := models.JpnCard{}
	if rows.Next() {
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

// GetRandCardForUser Вернет случайную карточку пользователя
func (r CardRepo) GetRandCardForUser(sessionId string) (models.JpnCard, error) {
	sql := "SELECT cards.* " +
		"FROM cards " +
		"LEFT JOIN sessions " +
		"	ON cards.userId = sessions.userId " +
		"WHERE sessions.sessionId = ? AND expires > NOW()" +
		"ORDER BY rand() LIMIT 1"

	rows, err := r.DB.Query(sql, sessionId)
	defer rows.Close()
	if err != nil {
		return models.JpnCard{}, err
	}

	card := models.JpnCard{}

	if rows.Next() {
		err = rows.Scan(&card.Id, &card.InJapan, &card.InRussian, &card.Mark, &card.DateAdd, &card.UserId)
		if err != nil {
			return models.JpnCard{}, err
		}
	}

	return card, nil
}

// GetCardByInJapan Вернет карточку по слову на японском
func (r CardRepo) GetCardByInJapan(inJapan string) (models.JpnCard, bool) {
	sql := "SELECT * FROM cards WHERE InJapan = ?"
	rows, err := r.DB.Query(sql, inJapan)
	defer rows.Close()
	if err != nil {
		return models.JpnCard{}, false
	}

	resultCard := models.JpnCard{}
	if rows.Next() {
		rows.Scan(
			&resultCard.Id,
			&resultCard.InJapan,
			&resultCard.InRussian,
			&resultCard.Mark,
			&resultCard.DateAdd,
			&resultCard.UserId)
	} else {
		return models.JpnCard{}, false
	}

	return resultCard, true
}

// UpdateCardMark Изменит значение оценки слова на mark
func (r CardRepo) UpdateCardMark(id, mark int) error {
	sql := "UPDATE cards SET mark = mark + ? WHERE id = ?"
	_, err := r.DB.Exec(sql, mark, id)
	if err != nil {
		return err
	}
	return nil
}
