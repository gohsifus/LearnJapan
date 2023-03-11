package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"time"
)

type SessionRepo struct {
	DB postgres.DB
}

func NewSessionRepo(db postgres.DB) *SessionRepo {
	return &SessionRepo{
		DB: db,
	}
}

// Add добавит сессию
func (r SessionRepo) Add(session *models.Session) error {
	query, args, err := squirrel.
		Insert("production.sessions").
		Columns("id", "user_id", "expires").
		Values(session.Id, session.UserId, session.Expires).
		ToSql()

	if err != nil {
		return err
	}

	return r.DB.Raw(query, args...).Scan(&session).Error
}

// DeleteSession удаляет сессию
func (r SessionRepo) DeleteSession(sessionId string) error {
	query, args, err := squirrel.
		Delete("production.sessions").
		Where(squirrel.Eq{"id": sessionId}).
		ToSql()

	if err != nil {
		return err
	}

	return r.DB.Raw(query, args...).Error
}

// UpdateSessionExpires обновит время жизни сессии
func (r SessionRepo) UpdateSessionExpires(sessionId, newExpires string) (bool, error) {
	return true, nil
}

// IsAliveSession проверит, действует ли сессия
func (r SessionRepo) IsAliveSession(sessionId string) (result bool, err error) {
	query, args, err := squirrel.Select("1").
		From("production.sessions").
		Where(squirrel.Eq{"id": sessionId}).
		Prefix("select exists(").
		Suffix(")").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// GetUserIdBySessionId вернет id пользователя по id сессии
func (r SessionRepo) GetUserIdBySessionId(sessionId string) (result int, err error) {
	query, args, err := squirrel.
		Select("user_id").
		From("production.sessions").
		Where(squirrel.Eq{"id": sessionId}).
		Where("expires > NOW()").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}

// Now вернет текущее время и дату на сервере с бд
func (r SessionRepo) Now() (result time.Time, err error) {
	query, args, err := squirrel.
		Select("NOW()").
		From("production.sessions").
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}
