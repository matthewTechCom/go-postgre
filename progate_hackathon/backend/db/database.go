package db

import (
	"database/sql"
	"fmt"
	"log"
	
	"github.com/matthewTechCom/progate_hackathon/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("PostgreSQL接続に失敗しました: %v", err)
	}
	// データベースが応答しているか
	if err := db.Ping(); err != nil {
		log.Fatalf("PostgreSQLへのPingに失敗しました: %v", err)
	}
	log.Println("PostgreSQL接続に成功しました")
	return db
}