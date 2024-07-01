package controller

import (
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/controller"
	"github.com/Jackzode/painting/service/question"
	"github.com/gin-gonic/gin"
)

// QuestionController question controller
type QuestionController struct {
	QuestionService *question.QuestionService
}

// NewQuestionController new controller
func NewQuestionController() *QuestionController {
	return &QuestionController{
		QuestionService: question.NewQuestionService(),
	}
}

func (qc *QuestionController) AddQuestion(ctx *gin.Context) {

	req := &types.QuestionAdd{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}

	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)

	//add question into db
	resp, err := qc.QuestionService.AddQuestion(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}

	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (qc *QuestionController) GetQuestion(ctx *gin.Context) {
	id := ctx.Query("id")
	id = utils.DeShortID(id)
	userID, _ := utils.GetUidFromTokenByCtx(ctx)
	//校对当前用户是否是obj的creator
	info, err := qc.QuestionService.GetQuestionAndAddPV(ctx, id, userID)
	if err != nil {
		return
	}
	info.ID = utils.EnShortID(info.ID)
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, info)
}

func (qc *QuestionController) GetQuestionPage(ctx *gin.Context) {
	req := &types.QuestionPageReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		return
	}
	req.LoginUserID, _ = utils.GetUidFromTokenByCtx(ctx)
	fmt.Printf("GetQuestionPage-req=%+v\n", req)
	questions, total, err := qc.QuestionService.GetQuestionPage(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, utils.NewPageModel(total, questions))
}

func (qc *QuestionController) GetQuestionInviteUserInfo(ctx *gin.Context) {
	questionID := utils.DeShortID(ctx.Query("id"))
	resp, err := qc.QuestionService.InviteUserInfo(ctx, questionID)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)

}

func (qc *QuestionController) UpdateQuestion(ctx *gin.Context) {
	req := &types.QuestionUpdate{}

	if !controller.BindAndCheckParams(ctx, req) {
		return
	}
	req.ID = utils.DeShortID(req.ID)
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	resp, err := qc.QuestionService.UpdateQuestion(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	respInfo, ok := resp.(*types.QuestionInfo)
	if !ok {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}

	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, &types.UpdateQuestionResp{UrlTitle: respInfo.UrlTitle, WaitForReview: !req.NoNeedReview})
}

func (qc *QuestionController) PersonalQuestionPage(ctx *gin.Context) {
	req := &types.PersonalQuestionPageReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.ParamErr, nil)
		return
	}

	req.LoginUserID, _ = utils.GetUidFromTokenByCtx(ctx)
	fmt.Printf("PersonalQuestionPage-req=%+v\n", req)
	resp, err := qc.QuestionService.PersonalQuestionPage(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)

}

func (qc *QuestionController) PersonalCollectionPage(ctx *gin.Context) {
	req := &types.PersonalCollectionPageReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		return
	}

	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	resp, err := qc.QuestionService.PersonalCollectionPage(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}
