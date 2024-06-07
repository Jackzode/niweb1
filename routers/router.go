package routers

import (
	"github.com/Jackzode/painting/middleware"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(debug bool) *gin.Engine {

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.GET("/heartBeat", heartBeats)
	r.Use(middleware.Cors())
	group := r.Group("/painting", middleware.InOutLog())
	InitUserRoutes(group)

	return r
}

func heartBeats(ctx *gin.Context) {

	ctx.String(200, "OK I am heartBeats")
}
