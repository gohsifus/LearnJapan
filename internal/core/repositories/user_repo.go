package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"database/sql"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// Add добавит пользователя в базу
func (r UserRepo) Add(user *models.User) error {
	sql := "INSERT INTO users (login, password, email) VALUES(?, ?, ?)"
	if _, err := r.DB.Exec(sql, user.Login, user.Password, user.Email); err != nil {
		return err
	}

	return nil
}

// FindUserByLoginAndPassword вернет пользователя если существует в базе
func (r UserRepo) FindUserByLoginAndPassword(login, password string) (models.User, bool) {
	sql := "SELECT id, login, email FROM users WHERE login = ? AND password = ?"
	rows, err := r.DB.Query(sql, login, password)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	user := models.User{}
	if rows.Next() {
		rows.Scan(&user.Id, &user.Login, &user.Email)
		return user, true
	}

	return user, false
}
