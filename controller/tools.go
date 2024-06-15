package controller

import (
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindAndCheckParams(ctx *gin.Context, data interface{}) bool {
	lang := utils.GetLang(ctx)
	ctx.Set(constants.AcceptLanguageFlag, lang)
	if err := ctx.ShouldBind(data); err != nil {
		glog.Slog.Errorf("http_handle BindAndCheckParams fail, %s", err.Error())
		HandleResponse(ctx, constants.ParamInvalid, err.Error(), nil)
		return false
	}
	return true
}

func HandleResponse(ctx *gin.Context, code int, msg string, data interface{}) {
	trace, _ := ctx.Get(constants.TraceID)
	bodyData := NewRespBodyData(code, msg, trace, data)
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, bodyData)
	return

}

func NewRespBodyData(code int, msg string, trace interface{}, data interface{}) *RespBody {
	return &RespBody{
		Code:    code,
		Message: msg,
		TraceId: trace,
		Data:    data,
	}
}

type RespBody struct {
	// http code
	Code int `json:"code"`
	// response message
	Message string `json:"msg"`
	//trace_id
	TraceId interface{} `json:"trace_id"`
	// response data
	Data interface{} `json:"data"`
}
