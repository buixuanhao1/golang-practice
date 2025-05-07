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

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error path parameter": err.Error(),
			})
			return
		}
		store := storage.NewSqlStore(db)
		business := biz.NewGetItemBiz(store)

		data, err := business.GetItemById(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error ": err.Error(),
			})
			return
		}
		// gin.H{
		// 	"data": data.Id,
		// }
		c.JSON(http.StatusAccepted, common.SimpleSuccessResponse(data))

	}
}
