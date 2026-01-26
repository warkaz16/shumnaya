package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(logger *slog.Logger) *gorm.DB {
	var err error

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "12345"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "shumnaya"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: dsn,
			// Используем подготовленные выражения и расширенный протокол для ускорения повторяющихся запросов
			PreferSimpleProtocol: false,
		}),
		&gorm.Config{PrepareStmt: true},
	)

	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	log.Println("БД успешно подключена")

	// Настройки пула соединений для повышения производительности
	sqlDB, err := db.DB()
	if err == nil {
		// Разумные дефолты; при необходимости переопределяйте через переменные окружения
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(0) // без лимита, управляется Postgres
	} else {
		logger.Warn("не удалось получить sql.DB для настройки пула", "error", err)
	}
	return db
}
