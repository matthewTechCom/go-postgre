package repository

import (
	"database/sql"
	"github.com/matthewTechCom/progate_hackathon/model"
)

type BoardRepositoryInterface interface {
	Save(board *model.Board) (int, error)
	GetByMiroID(miroBoardID string) (*model.Board, error)
}

type BoardRepository struct {
	DB *sql.DB
}

// コンストラクタ
func NewBoardRepository(db *sql.DB) BoardRepositoryInterface {
	return &BoardRepository{DB: db}
}

// ボード情報を保存する
func (r *BoardRepository) Save(board *model.Board) (int, error) {
	query := `
		INSERT INTO board (miro_board_id, title, description, created_at)
		VALUES ($1, $2, $3, DEFAULT)
		RETURNING id
	`

	var id int
	err := r.DB.QueryRow(query, board.MiroBoardID, board.Title, board.Description).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *BoardRepository) GetByMiroID(miroBoardID string) (*model.Board, error) {
	query := `
		SELECT id, miro_board_id, title, description, created_at
		FROM board
		WHERE miro_board_id = $1
	`
	row := r.DB.QueryRow(query, miroBoardID)
	var board model.Board
	err := row.Scan(&board.ID, &board.MiroBoardID, &board.Title, &board.Description, &board.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

