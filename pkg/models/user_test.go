package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Add(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil{
		panic(err)
	}

	DB = db
	defer db.Close()

	newUser := User{
		Email: "qwe@qwe.com",
		Login: "aboba",
		Password: "qweqweqwe",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(newUser.Login, newUser.Password, newUser.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	got, err := newUser.Add()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, got)
}

func TestFindUserByLoginAndPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil{
		panic(err)
	}

	DB = db
	defer db.Close()

	want := User{}

	type args struct{
		login string
		password string
	}

	newArgs := args{"test", "test"}

	mock.ExpectQuery("SELECT id, login, email FROM users").WithArgs(newArgs.login, newArgs.password).
		WillReturnRows(sqlmock.NewRows([]string{}))

	got, _ := FindUserByLoginAndPassword(newArgs.login, newArgs.password)

	assert.Equal(t, want, got)
}
