package routers

import (
	controller "github.com/Jackzode/painting/controller/user"
	"github.com/Jackzode/painting/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {

	g := r.Group("/user")
	userController := controller.NewUserController()
	g.POST("/registerByEmail", userController.UserRegisterByEmail)
	g.GET("/email/verification", userController.UserVerifyEmail)
	g.GET("/captcha", userController.UserRegisterCaptcha)
	g.POST("/login", userController.UserEmailLogin)
	g.GET("/getOtherProfile", userController.GetOtherUserInfoByUsername)

	//need login
	auth := g.Group("/personal", middleware.AccessToken())
	auth.GET("/profile", userController.GetUserInfoByUserID)
	auth.POST("/profile", userController.UpdateUserInfo)
	auth.POST("/uploadAvatar", userController.UploadAvatar)
}
