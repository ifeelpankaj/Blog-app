package db

import (
	"database/sql"
	"fmt"
	"time"

	"blog_app/config"
	"blog_app/utils/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Database struct {
	conn       *sql.DB
	isShutdown bool
}

var instance *Database

func Init() error {
	if instance != nil {
		return nil
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.DB_user,
		config.AppConfig.DB_pass,
		config.AppConfig.DB_host,
		config.AppConfig.DB_port,
		config.AppConfig.DB_name,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Error("Failed to open DB connection", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(30 * time.Second)

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping DB", zap.Error(err))
		return err
	}

	instance = &Database{conn: db}
	logger.Info("Database connection established", zap.String("DB_NAME", config.AppConfig.DB_name))
	return nil
}

func GetInstance() *Database {
	return instance
}

func (db *Database) HealthCheck() error {
	if db.isShutdown {
		return fmt.Errorf("database is shutting down")
	}
	if err := db.conn.Ping(); err != nil {
		logger.Error("Database connection is unhealthy", zap.Error(err))
		return err
	}
	logger.Info("Database connection is healthy", zap.String("Success", "True"))
	return nil
}

func (db *Database) Close() error {
	if db.isShutdown {
		return nil
	}
	db.isShutdown = true
	if err := db.conn.Close(); err != nil {
		logger.Error("Error closing database connection", zap.Error(err))
		return err
	}
	logger.Info("Database connection closed successfully")
	return nil
}
