package common

import "time"

type SQLModel struct {
	Id      int        `JSON:"id" gorm:"column:id;"`
	Created *time.Time `JSON:"created_at" gorm:"column:created_at;"`
	Updated *time.Time `JSON:"updated_at" gorm:"column:updated_at;"`
}
