package main

import (
	"fmt"
	"log"

	"github.com/matthewTechCom/progate_hackathon/chatgptapi"
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

	// テーブル作成
	migrate.Migrate(cfg)

	// リポジトリの初期化
	boardRepo := repository.NewBoardSummaryRepository(dbConn)

	// 外部APIクライアントの初期化
	miroClient := miroapi.NewMiroAPI(cfg.MiroAPIToken) 
	chatGPTClient := chatgptapi.NewChatGPTAPI(cfg.OpenAIApiKey)

	// ユースケースの初期化
	boardUsecase := usecase.NewBoardSummaryUsecase(boardRepo, miroClient, chatGPTClient)
	
	// Validatorの初期化
	validator := validator.NewValidator()

	// コントローラーの初期化
	boardController := controller.NewBoardSummaryController(boardUsecase, validator)

	// echoのルーター設定
	e := echo.New()

	// CORS ミドルウェアを適用
	e.Use(middleware.CORSMiddleware())
	// CSRF ミドルウェアを適用
	e.Use(middleware.CSRFMiddleware())

	router.SetupRoutes(e, boardController)

	// サーバー起動
	start := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("サーバースタート %s", start)
	if err := e.Start(start); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}