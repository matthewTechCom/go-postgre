package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser       string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
	ServerPort   string
	MiroAPIToken string
	DefaultBoardID  string
	DefaultAccessToken string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
        DBUser:       os.Getenv("DB_USER"),
        DBPassword:   os.Getenv("DB_PASSWORD"),
        DBHost:       os.Getenv("DB_HOST"),
        DBPort:       os.Getenv("DB_PORT"),
        DBName:       os.Getenv("DB_NAME"),
        ServerPort:   os.Getenv("SERVER_PORT"),
		MiroAPIToken: os.Getenv("MIRO_ACCESS_TOKEN"),
		DefaultBoardID:  os.Getenv("MIRO_BOARD_ID"), 
		DefaultAccessToken: os.Getenv("MIRO_ACCESS_TOKEN"), 
    }
}