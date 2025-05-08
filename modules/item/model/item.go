package model

import (
	"errors"
	"myginapp/common"
)

const (
	EntityName = "Item"
)

var (
	ErrTitleIsBlank = errors.New("Title can not be blank!")
	ErrItemDeleted  = errors.New("Item is deleted!")
)

type TodoItem struct {
	common.SQLModel
	UserId      int         `json:"user_id" gorm:"column:user_id"`
	Title       string      `JSON:"title" gorm:"column:title;"`
	Description string      `JSON:"description" gorm:"column:description;"`
	Status      *ItemStatus `JSON:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string { return "Todo_items" }

type TodoItemCreation struct {
	Id          int         `JSON:"-" gorm:"column:id;"`
	UserId      int         `json:"-" gorm:"column:user_id"`
	Title       string      `JSON:"title" gorm:"column:title;"`
	Description string      `JSON:"description" gorm:"column:description;"`
	Status      *ItemStatus `JSON:"status" gorm:"column:status;"`
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string     `JSON:"title" gorm:"column:title;"`
	Description *string     `JSON:"description" gorm:"column:description;"`
	Status      *ItemStatus `JSON:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }
