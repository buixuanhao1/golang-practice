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

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		var filter model.Filter
		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		paging.Process()

		store := storage.NewSqlStore(db)
		business := biz.NewListItemBiz(store)
		result, err := business.ListItem(c.Request.Context(), &filter, &paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error ": err.Error(),
			})
			return
		}
		// gin.H{
		// 	"data":   result,
		// 	"paging": paging,
		// }
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, &filter))

	}
}
