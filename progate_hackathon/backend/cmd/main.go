// cmd/main.go
package main

import (
	"fmt"
	"log"

	"github.com/matthewTechCom/progate_hackathon/config"
	"github.com/matthewTechCom/progate_hackathon/controller"
	"github.com/matthewTechCom/progate_hackathon/db"
	"github.com/matthewTechCom/progate_hackathon/migrate"
	"github.com/matthewTechCom/progate_hackathon/miroapi"
	"github.com/matthewTechCom/progate_hackathon/middleware"
	"github.com/matthewTechCom/progate_hackathon/repository"
	"github.com/matthewTechCom/progate_hackathon/router"
	"github.com/matthewTechCom/progate_hackathon/usecase"
	"github.com/matthewTechCom/progate_hackathon/validator"
	"github.com/labstack/echo/v4"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()

	// DB接続の初期化
	dbConn := db.InitDB(cfg)

	// マイグレーション：board, stickyテーブル作成
	migrate.Migrate(cfg)

	// リポジトリの初期化
	boardRepo := repository.NewBoardRepository(dbConn)
	stickyRepo := repository.NewStickyRepository(dbConn)

	// 外部APIクライアントの初期化（Miro API）
	miroClient := miroapi.NewMiroAPI(cfg.MiroAPIToken)

	// ユースケースの初期化（Miroから情報を取得し、DBに保存する処理）
	widgetUsecase := usecase.NewWidgetUsecase(boardRepo, stickyRepo, miroClient)

	// Validatorの初期化
	v := validator.NewValidator()

	// コントローラーの初期化
	widgetController := controller.NewWidgetController(widgetUsecase, v)

	// echoのルーター設定
	e := echo.New()
	e.Use(middleware.CORSMiddleware())
	router.SetupRoutes(e, widgetController)

	// サーバー起動
	start := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("サーバースタート %s", start)
	if err := e.Start(start); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}
