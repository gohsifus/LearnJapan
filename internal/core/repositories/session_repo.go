package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"database/sql"
	"time"
)

type SessionRepo struct {
	DB *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{
		DB: db,
	}
}

// Add добавит сессию
func (r SessionRepo) Add(session *models.Session) error {
	sql := "INSERT INTO sessions (sessionId, userId, expires) VALUES(?, ?, ?)"

	_, err := r.DB.Exec(sql, session.SessionId, session.UserId, session.Expires)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSession удаляет сессию
func (r SessionRepo) DeleteSession(sessionId string) error {
	sql := "DELETE FROM sessions WHERE sessionId = ?"
	_, err := r.DB.Exec(sql, sessionId)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSessionExpires обновит время жизни сессии
func (r SessionRepo) UpdateSessionExpires(sessionId, newExpires string) (bool, error) {
	return true, nil
}

// IsAliveSession проверит, действует ли сессия
func (r SessionRepo) IsAliveSession(sessionId string) (bool, error) {
	sql := "SELECT COUNT(sessionId) FROM sessions WHERE sessionId = ? AND expires > NOW();"

	rows, err := r.DB.Query(sql, sessionId)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	var exist int
	if rows.Next() {
		rows.Scan(&exist)
	}

	if exist == 1 {
		return true, nil
	}

	return false, nil
}

// GetUserIdBySessionId вернет id пользователя по id сессии
func (r SessionRepo) GetUserIdBySessionId(sessionId string) (int, bool) {
	sql := "SELECT userId FROM sessions WHERE sessionId = ? AND expires > NOW();"
	rows, err := r.DB.Query(sql, sessionId)
	defer rows.Close()
	if err != nil {
		return 0, false
	}

	var ret int
	if rows.Next() {
		rows.Scan(&ret)
		return ret, true
	} else {
		return 0, false
	}
}

// Now вернет текущее время и дату на сервере с бд
func (r SessionRepo) Now() (time.Time, error) {
	sql := "SELECT NOW()"

	rows, err := r.DB.Query(sql)
	defer rows.Close()

	if err != nil {
		return time.Time{}, err
	}

	var ret time.Time
	if rows.Next() {
		rows.Scan(&ret)
	}

	return ret, nil
}
