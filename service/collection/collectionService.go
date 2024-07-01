package collectService

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/collection"
	"github.com/jinzhu/copier"
)

type CollectionService struct {
	collection *collection.CollectionDao
	groupDao   *collection.CollectionGroupDao
}

func NewCollectionService() *CollectionService {
	return &CollectionService{
		collection: collection.NewCollectionDao(),
		groupDao:   collection.NewCollectionGroupRepo(),
	}
}

// AddCollectionGroup add collection group
func (cs *CollectionService) AddCollectionGroup(ctx context.Context, req *types.AddCollectionGroupReq) (err error) {
	collectionGroup := &types.CollectionGroup{}
	_ = copier.Copy(collectionGroup, req)
	return cs.groupDao.AddCollectionGroup(ctx, collectionGroup)
}

// UpdateCollectionGroup update collection group
func (cs *CollectionService) UpdateCollectionGroup(ctx context.Context, req *types.UpdateCollectionGroupReq, cols []string) (err error) {
	collectionGroup := &types.CollectionGroup{}
	_ = copier.Copy(collectionGroup, req)
	return cs.groupDao.UpdateCollectionGroup(ctx, collectionGroup, cols)
}

// GetCollectionGroup get collection group one
func (cs *CollectionService) GetCollectionGroup(ctx context.Context, id string) (resp *types.GetCollectionGroupResp, err error) {
	collectionGroup, exist, err := cs.groupDao.GetCollectionGroup(ctx, id)
	if err != nil {
		return
	}
	if !exist {
		return nil, errors.New(constants.UnknownError)
	}

	resp = &types.GetCollectionGroupResp{}
	_ = copier.Copy(resp, collectionGroup)
	return resp, nil
}

func (cs *CollectionService) SearchObjectCollectedByIds(ctx context.Context, userId string, objectIds []string) (collectedMap map[string]bool, err error) {
	return cs.collection.SearchObjectCollected(ctx, userId, objectIds)
}

func (cs *CollectionService) SearchList(ctx context.Context, search *types.CollectionSearch) ([]*types.Collection, int64, error) {
	return cs.collection.SearchList(ctx, search)
}

func (cs *CollectionService) CollectionSwitch(ctx context.Context, req *types.CollectionSwitchReq) (
	resp *types.CollectionSwitchResp, err error) {

	collectionGroup, err := cs.groupDao.CreateDefaultGroupIfNotExist(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	collect, exist, err := cs.collection.GetOneByObjectIDAndUser(ctx, req.UserID, req.ObjectID)
	if err != nil {
		return nil, err
	}
	if (!req.Bookmark && !exist) || (req.Bookmark && exist) {
		return nil, nil
	}

	if req.Bookmark {
		collect = &types.Collection{
			UserID:                req.UserID,
			ObjectID:              req.ObjectID,
			UserCollectionGroupID: collectionGroup.ID,
		}
		err = cs.collection.AddCollection(ctx, collect)
	} else {
		err = cs.collection.RemoveCollection(ctx, collect.ID)
	}
	if err != nil {
		return nil, err
	}

	// For now, we only support bookmark for question, so we just update question collect count
	resp = &types.CollectionSwitchResp{}
	//resp.ObjectCollectionCount, err = QuestionCommonServicer.UpdateCollectionCount(ctx, req.ObjectID)
	//if err != nil {
	//	return nil, err
	//}
	return resp, nil
}
