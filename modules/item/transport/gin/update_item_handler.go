package ginItem

import (
	"myginapp/common"
	"myginapp/modules/item/biz"
	"myginapp/modules/item/model"
	"myginapp/modules/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemUpdate
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := storage.NewSqlStore(db)
		business := biz.NewUpdateItemBiz(store, requester)
		if err := business.UpdateItemById(c.Request.Context(), id, &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// gin.H{
		// 	"data": true,
		// }
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}
