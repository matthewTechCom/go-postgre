// usecase/widget_usecase.go
package usecase

import (
	"fmt"
	"strings"

	"github.com/matthewTechCom/progate_hackathon/model"
	"github.com/matthewTechCom/progate_hackathon/miroapi"
	"github.com/matthewTechCom/progate_hackathon/repository"
)

type WidgetUsecaseInterface interface {
	ProcessAndSave(boardID, accessToken string) ([]int, error)
}

type WidgetUsecase struct {
	BoardRepo  repository.BoardRepositoryInterface
	StickyRepo repository.StickyRepositoryInterface
	MiroAPI    miroapi.MiroAPIInterface
}

func NewWidgetUsecase(boardRepo repository.BoardRepositoryInterface, stickyRepo repository.StickyRepositoryInterface, miro miroapi.MiroAPIInterface) WidgetUsecaseInterface {
	return &WidgetUsecase{
		BoardRepo:  boardRepo,
		StickyRepo: stickyRepo,
		MiroAPI:    miro,
	}
}

func (u *WidgetUsecase) ProcessAndSave(boardID, accessToken string) ([]int, error) {
	// DBに存在するかチェックし、なければ新規保存する
	board, err := u.BoardRepo.GetByMiroID(boardID)
	if err != nil {
		// 存在しない場合、新たな board を保存
		newBoard := &model.Board{
			MiroBoardID: boardID,
			Title:       "",
			Description: "",
		}
		boardIDInt, err := u.BoardRepo.Save(newBoard)
		if err != nil {
			return nil, fmt.Errorf("boardの保存に失敗: %v", err)
		}
		board = &model.Board{ID: boardIDInt, MiroBoardID: boardID}
	}

	// Miro APIからウィジェット情報を取得
	widgets, err := u.MiroAPI.GetWidgets(boardID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("miro APIから情報取得に失敗: %v", err)
	}

	var stickies []*model.Sticky
	for _, widget := range widgets {
		// ここでは widget.Text に「改善」または「反省」が含まれている場合のみ保存する
		var category string
		if strings.Contains(widget.Text, "改善") {
			category = "改善点"
		} else if strings.Contains(widget.Text, "反省") {
			category = "反省点"
		} else {
			// 該当しない場合はスキップ
			continue
		}

		sticky := &model.Sticky{
			BoardID:      board.ID,
			MiroStickyID: widget.ID,
			Content:      widget.Text,
			Category:     category,
		}

		stickies = append(stickies, sticky)
	}

	// 複数の付箋を保存
	savedIDs, err := u.StickyRepo.Save(stickies)
	if err != nil {
		return nil, fmt.Errorf("stickyの保存に失敗: %v", err)
	}

	return savedIDs, nil
}
