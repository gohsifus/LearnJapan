package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/postgres"
	"github.com/Masterminds/squirrel"
)

type CardRepo struct {
	DB postgres.DB
}

func NewCardRepo(db postgres.DB) *CardRepo {
	return &CardRepo{
		DB: db,
	}
}

// GetCardById Возвращает карточку по id
func (r CardRepo) GetCardById(id string) (result models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("*").
		From("production.cards").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// GetCardList Возвращает все карточки
func (r CardRepo) GetCardList() (result []models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("*").
		From("production.cards").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// GetCardListBySessionId Вернет карточки принадлежащие конкретному пользователю по id действующей сессии
func (r CardRepo) GetCardListBySessionId(sessionId string) (result []models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("cards.*").
		From("production.cards").
		LeftJoin("production.sessions ON cards.user_id = sessions.user_id").
		Where(squirrel.Eq{"sessions.id": sessionId}).
		Where("expires > NOW()").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// Add добавляет карточку
func (r CardRepo) Add(card *models.JpnCard) error {
	query, args, err := squirrel.
		Insert("production.cards").
		Columns("in_japan", "in_russian", "mark", "date_add", "user_id").
		Values(card.InJapan, card.InRussian, card.Mark, card.DateAdd, card.UserId).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return err
	}

	return r.DB.Raw(query, args...).Scan(&card).Error
}

// GetRandCard Вернет случайную карточку
func (r CardRepo) GetRandCard() (result models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("*").
		From("production.cards").
		OrderBy("random()").
		Limit(1).
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// GetRandCardForUser Вернет случайную карточку пользователя
func (r CardRepo) GetRandCardForUser(sessionId string) (result models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("cards.*").
		From("production.cards").
		LeftJoin("production.sessions ON cards.user_id = sessions.user_id").
		Where(squirrel.Eq{"sessions.id": sessionId}).
		Where("expires > NOW()").
		OrderBy("random()").
		Limit(1).
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// GetCardByInJapan Вернет карточку по слову на японском
func (r CardRepo) GetCardByInJapan(inJapan string) (result models.JpnCard, err error) {
	query, args, err := squirrel.
		Select("*").
		From("production.cards").
		Where(squirrel.Eq{"in_japan": inJapan}).
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// UpdateCardMark Изменит значение оценки слова на mark
func (r CardRepo) UpdateCardMark(id, mark int) (result models.JpnCard, err error) {
	query, args, err := squirrel.
		Update("production.cards").
		Set("mark", squirrel.Expr("mark + ?", mark)).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}
