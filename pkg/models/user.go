package models

func (user *User) Add() (bool, error){

	sql := "INSERT INTO users (login, password, email) VALUES(?, ?, ?)"
	if _, err := DB.Exec(sql, user.Login, user.Password, user.Email); err != nil{
		return false, err
	}

	return true, nil
}
