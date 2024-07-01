package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InOutLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//in log
		fmt.Println(ctx.FullPath())
		//next step
		ctx.Next()
		//fmt.Println(ctx.Request.RequestURI)
		//out log
	}
}
