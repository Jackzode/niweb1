package middleware

import (
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/gin-gonic/gin"
)

func TraceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := utils.GetTraceIdFromHeader(ctx)
		if traceId == "" {
			traceId = utils.GenerateTraceId()
		}
		ctx.Set(constants.TraceID, traceId)
		ctx.Next()
	}

}
