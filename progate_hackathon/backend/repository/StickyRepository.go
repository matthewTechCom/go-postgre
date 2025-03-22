package repository

import (
	"database/sql"
	"github.com/matthewTechCom/progate_hackathon/model"
)

type StickyRepositoryInterface interface {
	Save(sticky []*model.Sticky) ([]int, error)
}

type StickyRepository struct {
	DB *sql.DB
}

func NewStickyRepository(db *sql.DB) StickyRepositoryInterface {
	return &StickyRepository{DB: db}
}

func (r *StickyRepository) Save(stickies []*model.Sticky) ([]int, error) {
	var savedIDs []int

	for _, sticky := range stickies {
		query := `
			INSERT INTO sticky (board_id, miro_sticky_id, content, category, created_at)
			VALUES ($1, $2, $3, $4, DEFAULT)
            RETURNING id
		`
		var id int
		err := r.DB.QueryRow(query, sticky.BoardID, sticky.MiroStickyID, sticky.Content, sticky.Category).Scan(&id)
		if err != nil {
			return nil, err
		}
		savedIDs = append(savedIDs, id)
	}

	return savedIDs, nil

}
