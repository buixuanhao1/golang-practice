package biz

import (
	"context"
	"errors"
	"myginapp/common"
	"myginapp/modules/item/model"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, data *model.TodoItemUpdate) error
}

// store có thể là sqlStore, hoặc sau này bạn có thể tạo memoryStore, mockStore, apiStore, v.v.
type updateItemBiz struct {
	store     UpdateItemStorage
	requester common.Requester
}

func NewUpdateItemBiz(store UpdateItemStorage, requester common.Requester) *updateItemBiz {
	return &updateItemBiz{store: store, requester: requester}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {

	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		if err == common.RecordNotFound {
			return common.ErrCannotGetEntity(model.EntityName, err)
		}
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	if data.Status != nil && *data.Status == model.ItemStatusDeleted {
		return common.ErrCannotDeleteEntity(model.EntityName, model.ErrItemDeleted)
	}

	isOwner := biz.requester.GetUserId() == data.Id

	if !isOwner && !common.IsAdmin(biz.requester) {
		return common.ErrNoPermission(errors.New("no permission"))
	}

	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	return nil
}
