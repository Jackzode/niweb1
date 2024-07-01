package routers

import (
	"github.com/Jackzode/painting/controller/collection"
	"github.com/Jackzode/painting/controller/comment"
	"github.com/Jackzode/painting/controller/likes"
	controller "github.com/Jackzode/painting/controller/post"
	"github.com/Jackzode/painting/middleware"
	"github.com/gin-gonic/gin"
)

func InitPostRoutes(r *gin.RouterGroup) {

	g := r.Group("/post")
	postController := controller.NewQuestionController()

	g.GET("/getPost", postController.GetQuestion)
	g.GET("/getPostByPage", postController.GetQuestionPage)

	g.POST("/addPost", middleware.AccessToken(), postController.AddQuestion)
	g.GET("/getPersonalPost", middleware.AccessToken(), postController.PersonalQuestionPage)
	g.GET("/getPersonalCollectionPost", middleware.AccessToken(), postController.PersonalCollectionPage)

}

func InitLikePostRoutes(r *gin.RouterGroup) {
	g := r.Group("/post")
	likesController := likes.NewLikesController()
	g.POST("/likes", middleware.AccessToken(), likesController.AddOrCancelLikes)
	g.GET("/getLikesCount", likesController.CountLikes)
	g.GET("/checkLiked", middleware.AccessToken(), likesController.CheckLiked)
}

func InitCollectionRoutes(r *gin.RouterGroup) {
	g := r.Group("/collection")
	collectionController := collection.NewCollectionController()
	g.POST("/save", middleware.AccessToken(), collectionController.CollectionSwitch)
	g.GET("/checkSaved", middleware.AccessToken(), collectionController.CheckCollection)

}

func InitCommentRoutes(r *gin.RouterGroup) {
	g := r.Group("/comments")
	Controller := comment.NewCommentController()
	g.POST("/add", middleware.AccessToken(), Controller.AddComment)
	g.GET("/getCommentsPage", middleware.AccessToken(), Controller.GetCommentWithPage)

}
