package ginUser

import (
	"fmt"
	"myginapp/common"
	"myginapp/modules/user/biz"
	"myginapp/modules/user/model"
	"myginapp/modules/user/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// In ra nội dung biến data trước khi gọi Register
		fmt.Printf("Data before register: %+v\n", data)

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		biz := biz.NewRegisterBusiness(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Print("alooo")

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
