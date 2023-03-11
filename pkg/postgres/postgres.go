package postgres

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/pkg/gorm_logger"
	"LearnJapan.com/pkg/logger"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	cfg    *configs.Configs
	logger *logger.Logger
}

func NewDBPostgres(cfg *configs.Configs, logger *logger.Logger) (db DB, err error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.DB,
	)

	gormDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: gorm_logger.NewGormLogger(cfg),
	})
	if err != nil {
		return db, err
	}

	dbInstance, err := gormDB.DB()
	if err != nil {
		return db, err
	}

	if err := dbInstance.Ping(); err != nil {
		return db, err
	}

	logger.Info("database connection established")

	return DB{
		DB:     gormDB,
		cfg:    cfg,
		logger: logger,
	}, nil
}
