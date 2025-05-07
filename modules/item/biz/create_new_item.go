package biz

import (
	"context"
	"myginapp/common"
	"myginapp/modules/item/model"
	"strings"
)

type CreateItemStorage interface {
	InsertItem(ctx context.Context, data *model.TodoItemCreation) error
}

// store có thể là sqlStore, hoặc sau này bạn có thể tạo memoryStore, mockStore, apiStore, v.v.
type createItemBiz struct {
	store CreateItemStorage
}

func NewCreateItemBiz(store CreateItemStorage) *createItemBiz {
	return &createItemBiz{store: store}
}

func (biz *createItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	title := strings.TrimSpace(data.Title)

	if title == "" {
		return model.ErrTitleIsBlank
	}

	if err := biz.store.InsertItem(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
