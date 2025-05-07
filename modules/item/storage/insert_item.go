package storage

import (
	"context"
	"myginapp/common"
	"myginapp/modules/item/model"
)

// hàm để phục vụ cho tầng business
func (s *sqlStore) InsertItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
