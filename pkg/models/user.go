package models

//Add Добавит пользователя в базу
func (user *User) Add() (bool, error){
	sql := "INSERT INTO users (login, password, email) VALUES(?, ?, ?)"
	if _, err := DB.Exec(sql, user.Login, user.Password, user.Email); err != nil{
		return false, err
	}

	return true, nil
}

//FindUserByLoginAndPassword Вернет пользователя если существует в базе
func FindUserByLoginAndPassword(login, password  string) (User, bool){
	sql := "SELECT id, login, email FROM users WHERE login = ? AND password = ?"
	rows, err := DB.Query(sql, login, password)
	defer rows.Close()
	if err != nil{
		panic(err)
	}

	user := User{}
	if rows.Next(){
		rows.Scan(&user.Id, &user.Login, &user.Email)
		return user, true
	}

	return user, false
}
