package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemStatus int

const (
	itemStatusDoing ItemStatus = iota
	itemStatusDone
	itemStatusDeleted
)

var allItemStatuses = [3]string{"Doing", "Done", "Deleted"}

func (item *ItemStatus) String() string {
	return allItemStatuses[*item]
}

func parseStringToItemStatus(s string) (ItemStatus, error) {
	for i := range allItemStatuses {
		if allItemStatuses[i] == s {
			return ItemStatus(i), nil
		}
	}
	return ItemStatus(0), errors.New("Invalid status string")
}

func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("fail to scan data from sql: %s", value)
	}
	v, err := parseStringToItemStatus(string(bytes))

	if err != nil {
		return fmt.Errorf("fail to scan data from sql %s", value)
	}
	*item = v

	return nil
}

func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

func (item *ItemStatus) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}
	return item.String(), nil
}

func (item *ItemStatus) UnmarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "")

	itemValue, err := parseStringToItemStatus(str)
	if err != nil {
		return err
	}

	*item = itemValue

	return nil
}

type TodoItem struct {
	Id          int         `JSON:"id" gorm:"column:id;"`
	Title       string      `JSON:"title" gorm:"column:title;"`
	Description string      `JSON:"description" gorm:"column:description;"`
	Status      *ItemStatus `JSON:"status" gorm:"column:status;"`
	Created     *time.Time  `JSON:"create_at" gorm:"column:create_at;"`
	Updated     *time.Time  `JSON:"update_at, omitempty" gorm:"column:update_at;" `
}

func (TodoItem) TableName() string { return "Todo_items" }

type TodoItemCreation struct {
	Id          int    `JSON:"-" gorm:"column:id;"`
	Title       string `JSON:"title" gorm:"column:title;"`
	Description string `JSON:"description" gorm:"column:description;"`
	Status      string `JSON:"status" gorm:"column:status;"`
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string `JSON:"title" gorm:"column:title;"`
	Description *string `JSON:"description" gorm:"column:description;"`
	Status      *string `JSON:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }

type Paging struct {
	Page  int   `JSON:"page" form:"page"`
	Limit int   `JSON:"limit" form:"limit"`
	Total int64 `JSON:"total" form:"total"`
}

func (Paging) TableName() string { return TodoItem{}.TableName() }

func (p *Paging) Process() {
	if p.Page < 0 {
		p.Page = 1
	}

	if p.Limit < 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

func main() {

	dsn := "root:mysecret@tcp(127.0.0.1:3308)/Todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PUT("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	r.Run(":3000")

}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error should bind": err.Error(),
			})
			return
		}

		paging.Process()
		var result []TodoItem

		if err := db.Table(TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Order("id asc").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error fetch by id": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})

	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItem
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error path parameter": err.Error(),
			})
			return
		}

		data.Id = id
		if err := db.First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error fetch by id": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})

	}
}

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error path parameter": err.Error(),
			})
			return
		}
		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
			"status": "Deleted",
		}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error delete by id": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})

	}
}

func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemCreation
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error should bind": err.Error(),
			})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error create": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"data": data.Id,
		})

	}
}

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemUpdate
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error path parameter": err.Error(),
			})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error should bind": err.Error(),
			})
			return
		}

		// Cập nhật bản ghi
		if err := db.Model(&TodoItem{}).Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})

	}
}
