package middleware

import (
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func AccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get token
		// 我们这里jwt鉴权取头部信息 x-token
		//登录时回返回token信息 这里前端需要把token
		//存储到cookie或者本地localStorage中
		//不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := ctx.Request.Header.Get(constants.HeaderToken)
		if token == "" {
			fmt.Println(token, "...token empty")
			//response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			ctx.Abort()
			return
		}
		// parseToken 解析token包含的信息
		claims, err := utils.ParseToken(token)
		if err != nil {
			ctx.Abort()
			return
		}
		//refresh token
		if claims.ExpiresAt.Unix()-time.Now().Unix() <= constants.ExpireBuffer {

		}
		ctx.Set(constants.TokenClaim, claims)
		ctx.Next()
	}
}
