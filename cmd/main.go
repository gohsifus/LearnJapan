package main

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/internal/core/repositories"
	"LearnJapan.com/internal/delivery/controllers"
	"LearnJapan.com/internal/delivery/middlewares"
	v1 "LearnJapan.com/internal/delivery/router/v1"
	"LearnJapan.com/pkg/logger"
	pg "LearnJapan.com/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	db, err := pg.NewDBPostgres(cfg, logs)
	if err != nil {
		logs.Fatal(err)
	}

	dbInstance, err := db.DB.DB()
	if err != nil {
		logs.Fatal(err)
	}

	driver, err := postgres.WithInstance(dbInstance, &postgres.Config{})
	if err != nil {
		logs.Fatal(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations/pg", "postgres", driver)
	if err != nil {
		logs.Fatal(err)
	}

	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		logs.Fatal(err)
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		logs.Info("Application is stopped")

		dbInstance.Close()

		os.Exit(0)
	}()

	sessionRepo := repositories.NewSessionRepo(db)
	authMiddleware := middlewares.NewAuthMiddleware(sessionRepo)

	mux := gin.New()
	controller := controllers.NewMainController(db, logs)
	router := v1.NewRouter(mux, controller, authMiddleware)
	router.Setup()

	if err := http.ListenAndServe(":8080", router.Mux); err != nil {
		logs.Fatal(err)
	}
}
