package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/frederickmarvel/supernova/internal/config"
	_ "github.com/lib/pq"
)

func New(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("db open error %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error %v", err)
	}
	return db
}
