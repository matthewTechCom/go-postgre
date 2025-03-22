package migrate

import (
	"database/sql"
	"log"

	"github.com/matthewTechCom/progate_hackathon/config"
	"github.com/matthewTechCom/progate_hackathon/db"
)

// Migrate データベースのマイグレーションを実行する
func Migrate(cfg *config.Config) {
	// db.goのInitDB関数を呼び出して、DB接続を取得
	dbConn := db.InitDB(cfg)

	// マイグレーション用のSQLを実行
	err := createBoardTable(dbConn)
	if err != nil {
		log.Fatalf("テーブル作成に失敗しました: %v", err)
	} else {
		log.Println("board テーブル作成が完了しました")
	}

	err = createStickyTable(dbConn)
	if err != nil {
		log.Fatalf("テーブル作成に失敗しました: %v", err)
	} else {
		log.Println("sticky テーブル作成が完了しました")
	}
}

// createBoardTable board テーブルを作成するSQL
func createBoardTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS board (
			id SERIAL PRIMARY KEY,
			miro_board_id TEXT NOT NULL,
			title TEXT,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// createStickyTable sticky テーブルを作成するSQL
func createStickyTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS sticky (
			id SERIAL PRIMARY KEY,
			board_id INT NOT NULL,
			miro_sticky_id TEXT NOT NULL,
			content TEXT NOT NULL,
			category TEXT CHECK (category IN ('改善点', '反省点')) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (board_id) REFERENCES board(id) ON DELETE CASCADE
		);
	`
	_, err := db.Exec(query)
	return err
}
