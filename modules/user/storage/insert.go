package storage

import (
	"context"
	"myginapp/common"
	"myginapp/modules/user/model"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *model.UserCreate) error {
	if err := s.db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
