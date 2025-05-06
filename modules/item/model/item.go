package model

import (
	"errors"
	"myginapp/common"
)

var (
	ErrTitleIsBlank = errors.New("Title can not be blank!")
)

type TodoItem struct {
	common.SQLModel
	Title       string      `JSON:"title" gorm:"column:title;"`
	Description string      `JSON:"description" gorm:"column:description;"`
	Status      *ItemStatus `JSON:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string { return "Todo_items" }

type TodoItemCreation struct {
	Id          int         `JSON:"-" gorm:"column:id;"`
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
