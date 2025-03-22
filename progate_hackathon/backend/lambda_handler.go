// lambda_handler.go
package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/matthewTechCom/progate_hackathon/config"
	"github.com/matthewTechCom/progate_hackathon/db"
	"github.com/matthewTechCom/progate_hackathon/migrate"
	"github.com/matthewTechCom/progate_hackathon/miroapi"
	"github.com/matthewTechCom/progate_hackathon/repository"
	"github.com/matthewTechCom/progate_hackathon/usecase"
)

type LambdaResponse struct {
	Message string
	IDs     []int
}

func HandleRequest(ctx context.Context) (LambdaResponse, error) {
	// 環境変数や設定ファイルから設定を読み込む
	cfg := config.LoadConfig()

	// DB接続の初期化
	dbConn := db.InitDB(cfg)

	// マイグレーション
	migrate.Migrate(cfg)

	// リポジトリの初期化
	boardRepo := repository.NewBoardRepository(dbConn)
	stickyRepo := repository.NewStickyRepository(dbConn)

	// Miro API クライアントの初期化
	miroClient := miroapi.NewMiroAPI(cfg.MiroAPIToken)

	// ユースケースの初期化
	widgetUsecase := usecase.NewWidgetUsecase(boardRepo, stickyRepo, miroClient)

	// BoardID、AccessToken は環境変数や Secrets Manager などで管理したい
	boardID := cfg.DefaultBoardID
	accessToken := cfg.DefaultAccessToken

	
	savedIDs, err := widgetUsecase.ProcessAndSave(boardID, accessToken)
	if err != nil {
		log.Printf("処理失敗: %v", err)
		return LambdaResponse{}, err
	}

	return LambdaResponse{
		Message: "情報をDBに保存しました",
		IDs:     savedIDs,
	}, nil
}

func main() {
	// Lambda ハンドラーとして起動
	lambda.Start(HandleRequest)
}
