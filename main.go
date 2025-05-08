package main

import (
	"log"
	"myginapp/component/tokenprovider/jwt"
	"myginapp/middleware"
	ginItem "myginapp/modules/item/transport/gin"
	"myginapp/modules/user/storage"
	ginUser "myginapp/modules/user/transport/gin"

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
	tokenProvider := jwt.NewTokenJWTProvider("jwt", "200Lab.io")
	authStore := storage.NewSQLStore(db)
	middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)
	r := gin.Default()
	r.Use(middleware.Recover())
	v1 := r.Group("/v1")
	{
		v1.POST("/register", ginUser.Register(db))
		v1.POST("/login", ginUser.Login(db, tokenProvider))
		v1.GET("/profile", middlewareAuth, ginUser.Profile())

		items := v1.Group("/items")
		{
			items.POST("", middlewareAuth, ginItem.CreateItem(db))
			items.GET("", ginItem.ListItem(db))
			items.GET("/:id", ginItem.GetItem(db))
			items.PUT("/:id", middlewareAuth, ginItem.UpdateItem(db))
			items.DELETE("/:id", middlewareAuth, ginItem.DeleteItem(db))
		}
	}

	r.Run(":3000")

}
