package biz

import (
	"context"
	"myginapp/common"
	"myginapp/modules/item/model"
)

type DeleteItemStorage interface {
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
}

// store có thể là sqlStore, hoặc sau này bạn có thể tạo memoryStore, mockStore, apiStore, v.v.
type deleteItemBiz struct {
	store DeleteItemStorage
}

func NewDeleteItemBiz(store DeleteItemStorage) *deleteItemBiz {
	return &deleteItemBiz{store: store}
}

func (biz *deleteItemBiz) DeleteItemById(ctx context.Context, id int) error {

	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		if err == common.RecordNotFound {
			return common.ErrCannotGetEntity(model.EntityName, err)
		}
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	if data != nil && *data.Status == model.ItemStatusDeleted {
		return common.ErrCannotDeleteEntity(model.EntityName, model.ErrItemDeleted)
	}

	err = biz.store.DeleteItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
