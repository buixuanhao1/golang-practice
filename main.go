package main

import (
	"log"
	"myginapp/common"
	"myginapp/modules/item/model"
	ginItem "myginapp/modules/item/transport/gin"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
			items.POST("", ginItem.CreateItem(db))
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
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error should bind": err.Error(),
			})
			return
		}

		paging.Process()
		var result []model.TodoItem

		if err := db.Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
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
		// gin.H{
		// 	"data":   result,
		// 	"paging": paging,
		// }
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))

	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItem
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
		// gin.H{
		// 	"data": data,
		// }
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))

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
		if err := db.Table(model.TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
			"status": "Deleted",
		}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error delete by id": err.Error(),
			})
			return
		}
		// gin.H{
		// 	"data": true,
		// }
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemUpdate
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
		if err := db.Model(&model.TodoItem{}).Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update",
			})
			return
		}

		// gin.H{
		// 	"data": true,
		// }
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}
