package ginItem

import (
	"myginapp/common"
	"myginapp/modules/item/biz"
	"myginapp/modules/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		store := storage.NewSqlStore(db)
		business := biz.NewDeleteItemBiz(store)
		if err := business.DeleteItemById(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// gin.H{
		// 	"data": true,
		// }
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}
