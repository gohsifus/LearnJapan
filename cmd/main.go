package main

import (
	_ "LearnJapan.com/cmd/router"
	"LearnJapan.com/pkg/configs"
	_ "LearnJapan.com/pkg/configs"
	"LearnJapan.com/pkg/logger"
	"LearnJapan.com/pkg/models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init(){
	db, err := sql.Open("mysql", configs.Cfg.DBConnectionString)
	if err != nil{
		panic("Ошибка подключения к базе")
	}
	
	dbStatus := db.Ping()
	if dbStatus != nil{
		fmt.Println("err: ")
		fmt.Println(dbStatus)
		logger.Print("err: " + dbStatus.Error())
	} else {
		fmt.Println("db connected")
		logger.Print("db connected")
	}
	
	models.DB = db
}

func main(){
	fmt.Println("Сервер запущен")
	logErr := logger.Print("Сервер запущен")

	if logErr != nil{
		panic(logErr)
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		fmt.Println("Сервер остановлен")
		logger.Print("Сервер остановлен")

		models.DB.Close()

		os.Exit(0)
	}()


	http.ListenAndServe("0.0.0.0:8080", nil)
}
