package repositories

import (
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/postgres"
	"github.com/Masterminds/squirrel"
)

type UserRepo struct {
	postgres.DB
}

func NewUserRepo(db postgres.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// Add добавит пользователя в базу
func (r UserRepo) Add(user *models.User) error {
	query, args, err := squirrel.
		Insert("production.users").
		Columns("login", "password", "email").
		Values(user.Login, user.Password, user.Email).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return err
	}

	return r.DB.Raw(query, args...).Scan(&user).Error
}

// FindUserByLoginAndPassword вернет пользователя если существует в базе
func (r UserRepo) FindUserByLoginAndPassword(login, password string) (result models.User, err error) {
	query, args, err := squirrel.
		Select("*").
		From("production.users").
		Where(squirrel.Eq{"login": login}).
		Where(squirrel.Eq{"password": password}).
		ToSql()

	if err != nil {
		return result, err
	}

	return result, r.DB.Raw(query, args...).Scan(&result).Error
}
