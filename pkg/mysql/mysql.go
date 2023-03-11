package mysql

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/pkg/logger"
	"database/sql"
	"fmt"
)

type DBMySql struct {
	Db     *sql.DB
	cfg    *configs.Configs
	logger *logger.Logger
}

func NewDBMySql(cfg *configs.Configs, logger *logger.Logger) (*DBMySql, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MYSQL.User,
		cfg.MYSQL.Password,
		cfg.MYSQL.Host,
		cfg.MYSQL.Port,
		cfg.MYSQL.DB,
	)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("database connection established")

	return &DBMySql{
		Db:     db,
		cfg:    cfg,
		logger: logger,
	}, nil
}
