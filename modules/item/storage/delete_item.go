package storage

import (
	"context"
	"myginapp/modules/item/model"
)

// hàm để phục vụ cho tầng business
func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	statusDeleted := model.ItemStatusDeleted
	if err := s.db.Table(model.TodoItem{}.TableName()).Where(cond).Updates(map[string]interface{}{
		"status": statusDeleted.String(),
	}).Error; err != nil {
		return err
	}
	return nil
}
