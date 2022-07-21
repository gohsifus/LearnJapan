package logger

import (
	"os"
	"time"
)

var logFile string

func init(){
	//TODO Путь можно брать из env
	logFile = "./logs/log.txt"
}

func Print(str string) error{
	file, err := os.OpenFile("./logs/log.txt", os.O_RDWR | os.O_APPEND, 0666)
	if err != nil{
		return err
	}
	defer file.Close()

	_, err = file.WriteString(time.Now().Format("02.01.2006 15:04:05 === ") + str + "\n")
	if err != nil{
		return err
	}

	return nil
}
