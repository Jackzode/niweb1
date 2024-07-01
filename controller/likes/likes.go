package likes

import (
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/controller"
	"github.com/Jackzode/painting/service/likes"
	"github.com/gin-gonic/gin"
)

type LikesController struct {
	like *likes.LikeService
}

func NewLikesController() *LikesController {
	return &LikesController{
		like: likes.NewLikeService(),
	}
}

func (lc *LikesController) AddOrCancelLikes(ctx *gin.Context) {
	req := &types.AddLikeReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}
	req.UserId, _ = utils.GetUidFromTokenByCtx(ctx)
	//fmt.Printf("req=%+v\n", req)
	var err error
	if req.Status {
		err = lc.like.AddLike(ctx, req)
	} else {
		err = lc.like.CancelLike(ctx, req)
	}
	if err != nil {
		glog.Slog.Error(err.Error())
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (lc *LikesController) CountLikes(ctx *gin.Context) {
	req := &types.AddLikeReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}
	req.UserId, _ = utils.GetUidFromTokenByCtx(ctx)
	//fmt.Printf("req=%+v\n", req)
	countLikes, err := lc.like.CountLikes(ctx, req)
	if err != nil {
		glog.Slog.Error(err.Error())
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	resp := map[string]interface{}{
		"count": countLikes,
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (lc *LikesController) CheckLiked(ctx *gin.Context) {
	req := &types.AddLikeReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}
	req.UserId, _ = utils.GetUidFromTokenByCtx(ctx)
	//fmt.Printf("req=%+v\n", req)
	liked, err := lc.like.CheckLiked(ctx, req)
	if err != nil {
		glog.Slog.Error(err.Error())
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	resp := map[string]interface{}{
		"liked": liked,
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}
