package models

import (
	"fmt"
	"time"
)

//Add Добавит сессию
func (s *Session) Add() (bool, error){
	sql := "INSERT INTO sessions (sessionId, userId, expires) VALUES(?, ?, ?)"
	_, err := DB.Exec(sql, s.SessionId, s.UserId, s.Expires)
	if err != nil{
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

//DeleteSession Удаляет сессию
func DeleteSession(sessionId string) (bool, error){
	sql := "DELETE FROM sessions WHERE sessionId = ?"
	_, err := DB.Exec(sql, sessionId)
	if err != nil {
		return false, err
	}

	return true, nil
}

//UpdateSessionExpires Обновит время жизни сессии
func UpdateSessionExpires(sessionId, newExpires string) (bool, error){
	return true, nil
}

//IsAliveSession Проверит действует ли сессия
func IsAliveSession(sessionId string) (bool, error){
	sql := "SELECT COUNT(sessionId) FROM sessions WHERE sessionId = ? AND expires > NOW();"

	rows, err := DB.Query(sql, sessionId)
	if err != nil{
		return false, err
	}

	var exist int
	if rows.Next(){
		rows.Scan(&exist)
	}

	if exist == 1 {
		return true, nil
	}

	return false, nil
}

//GetUserIdBySessionId Вернет id пользователя по id сессии
func GetUserIdBySessionId(sessionId string) (int, bool){
	sql := "SELECT userId FROM sessions WHERE sessionId = ? AND expires > NOW();"
	rows, err := DB.Query(sql, sessionId)
	if err != nil{
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

//Now Вернет текущее время и дату на сервере с бд
func Now() time.Time{
	sql := "SELECT NOW()"

	rows, err := DB.Query(sql)
	if err != nil{
		fmt.Println(err)
		return time.Time{}
	}

	var ret time.Time
	if rows.Next(){
		rows.Scan(&ret)
	}

	return ret
}
