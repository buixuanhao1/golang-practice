package storage

import (
	"context"
	"errors"
	"myginapp/common"
	"myginapp/modules/item/model"

	"gorm.io/gorm"
)

// hàm để phục vụ cho tầng business
func (s *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem
	if err := s.db.Where(cond).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
