package routes

import (
	"github.com/gin-gonic/gin"
	"micro-gin/app/controllers/app"
	"micro-gin/app/controllers/common"
	"micro-gin/app/middleware"
	"micro-gin/app/services"
)

func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.POST("/auth/register", app.Register)
	router.POST("/auth/test_gen_struct", app.TestGenStruct)
	router.POST("/auth/login", app.Login)

	authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))
	{
		authRouter.GET("/auth/info", app.Info)
		authRouter.POST("/auth/logout", app.Logout)
		authRouter.POST("/image_upload", common.ImageUpload)
	}
}
