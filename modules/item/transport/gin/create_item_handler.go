package ginItem

import (
	"myginapp/common"
	"myginapp/modules/item/biz"
	"myginapp/modules/item/model"
	"myginapp/modules/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemCreation
		// nhiệm vụ của handler
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		store := storage.NewSqlStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// gin.H{
		// 	"data": data.Id,
		// }
		c.JSON(http.StatusAccepted, common.SimpleSuccessResponse(data.Id))

	}
}
