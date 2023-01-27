package main

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/internal/core/repositories"
	"LearnJapan.com/internal/delivery/controllers"
	"LearnJapan.com/internal/delivery/middlewares"
	v1 "LearnJapan.com/internal/delivery/router/v1"
	"LearnJapan.com/pkg/logger"
	"LearnJapan.com/pkg/mysql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := configs.NewConfigs()
	if err != nil {
		log.Fatal(err)
	}

	logs := logger.NewLogger(cfg)

	db, err := mysql.NewDBMySql(cfg, logs)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		logs.Info("Application is stopped")

		db.Db.Close()

		os.Exit(0)
	}()

	sessionRepo := repositories.NewSessionRepo(db.Db)
	authMiddleware := middlewares.NewAuthMiddleware(sessionRepo)

	mux := gin.New()
	controller := controllers.NewMainController(db.Db, logs)
	router := v1.NewRouter(mux, controller, authMiddleware)
	router.Setup()

	if err := http.ListenAndServe(":8080", router.Mux); err != nil {
		logs.Fatal(err)
	}
}
