package storage

import (
	"context"
	"myginapp/modules/item/model"
)

// hàm để phục vụ cho tầng business
func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{},
	data *model.TodoItemUpdate) error {

	if err := s.db.Model(&model.TodoItem{}).Where(cond).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
