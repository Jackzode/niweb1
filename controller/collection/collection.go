package collection

import (
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/controller"
	"github.com/Jackzode/painting/service/collection"
	"github.com/gin-gonic/gin"
)

type CollectionController struct {
	collectionService *collectService.CollectionService
}

func NewCollectionController() *CollectionController {
	return &CollectionController{
		collectionService: collectService.NewCollectionService(),
	}
}

func (cc *CollectionController) CollectionSwitch(ctx *gin.Context) {
	req := &types.CollectionSwitchReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		return
	}

	req.ObjectID = utils.DeShortID(req.ObjectID)
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)

	resp, err := cc.collectionService.CollectionSwitch(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, err.Error(), nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (cc *CollectionController) CheckCollection(ctx *gin.Context) {
	req := &types.CollectionSwitchReq{}
	if !controller.BindAndCheckParams(ctx, req) {
		return
	}

	req.ObjectID = utils.DeShortID(req.ObjectID)
	req.UserID, _ = utils.GetUidFromTokenByCtx(ctx)
	collectedMap, err := cc.collectionService.SearchObjectCollectedByIds(ctx, req.UserID, []string{req.ObjectID})
	if err != nil {
		glog.Slog.Error(err.Error())
		controller.HandleResponse(ctx, constants.InternalErrCode, err.Error(), nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, collectedMap)
}
