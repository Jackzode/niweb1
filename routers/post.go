package routers

import (
	controller "github.com/Jackzode/painting/controller/post"
	"github.com/gin-gonic/gin"
)

func InitPostRoutes(r *gin.RouterGroup) {

	g := r.Group("/post")
	postController := controller.NewQuestionController()

	g.GET("/post", postController.GetQuestion)

}
