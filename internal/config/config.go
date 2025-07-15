package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int
}

func Load() *Config {
	godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
	}
	return &Config{
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
	}
}
