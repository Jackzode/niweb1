package comment

import (
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/controller"
	"github.com/Jackzode/painting/service/comment"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	cs *comment.CommentService
}

func NewCommentController() *CommentController {
	return &CommentController{}
}

func (cc *CommentController) AddComment(ctx *gin.Context) {
	req := &types.AddCommentReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}
	req.ObjectID = utils.DeShortID(req.ObjectID)
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	resp, err := cc.cs.AddComment(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (cc *CommentController) RemoveComment(ctx *gin.Context) {
	req := &types.RemoveCommentReq{}
	if controller.BindAndCheckParams(ctx, req) {
		return
	}
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	err := cc.cs.RemoveComment(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (cc *CommentController) GetCommentWithPage(ctx *gin.Context) {
	req := &types.GetCommentWithPageReq{}
	if controller.BindAndCheckParams(ctx, req) {
		return
	}
	req.ObjectID = utils.DeShortID(req.ObjectID)
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)

	resp, _, err := cc.cs.GetCommentWithPage(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (cc *CommentController) GetCommentPersonalWithPage(ctx *gin.Context) {
	req := &types.GetCommentPersonalWithPageReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		return
	}

	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)

	resp, _, err := cc.cs.GetCommentPersonalWithPage(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (cc *CommentController) GetComment(ctx *gin.Context) {
	req := &types.GetCommentReq{}
	if controller.BindAndCheckParams(ctx, req) {
		return
	}
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	resp, err := cc.cs.GetComment(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}
