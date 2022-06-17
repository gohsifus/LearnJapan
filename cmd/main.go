package main

import (
	_ "LearnJapan.com/cmd/router"
	"LearnJapan.com/pkg/configs"
	_ "LearnJapan.com/pkg/configs"
	"LearnJapan.com/pkg/models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	} else {
		fmt.Println("db connected")
	}
	
	models.DB = db

	fileLog, err := os.OpenFile("./logs/log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log.SetOutput(fileLog)
}

func main(){
	fmt.Println("Сервер запущен")
	log.Println("Сервер запущен")

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		fmt.Println("Сервер остановлен")
		log.Println("Сервер остановлен")

		models.DB.Close()

		os.Exit(0)
	}()

	http.ListenAndServe("0.0.0.0:8080", nil)
}
