package collection

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/tools"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

type CollectionGroupDao struct {
	DB    *xorm.Engine
	Cache *redis.Client
}

func NewCollectionGroupRepo() *CollectionGroupDao {
	return &CollectionGroupDao{
		DB:    handler.Engine,
		Cache: handler.RedisClient,
	}
}

// AddCollectionGroup add collection group
func (cr *CollectionGroupDao) AddCollectionGroup(ctx context.Context, collectionGroup *types.CollectionGroup) (err error) {
	_, err = cr.DB.Context(ctx).Insert(collectionGroup)
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

// AddCollectionDefaultGroup add collection group
func (cr *CollectionGroupDao) AddCollectionDefaultGroup(ctx context.Context, userID string) (collectionGroup *types.CollectionGroup, err error) {
	defaultGroup := &types.CollectionGroup{
		Name:         "default",
		DefaultGroup: constants.CGDefault,
		UserID:       userID,
	}
	_, err = cr.DB.Context(ctx).Insert(defaultGroup)
	if err != nil {
		err = errors.New(constants.DatabaseError)
		return
	}
	collectionGroup = defaultGroup
	return
}

// CreateDefaultGroupIfNotExist create default group if not exist
func (cr *CollectionGroupDao) CreateDefaultGroupIfNotExist(ctx context.Context, userID string) (
	collectionGroup *types.CollectionGroup, err error) {
	_, err = cr.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		old := &types.CollectionGroup{
			UserID:       userID,
			DefaultGroup: constants.CGDefault,
		}
		exist, err := session.ForUpdate().Get(old)
		if err != nil {
			return nil, err
		}
		if exist {
			collectionGroup = old
			return old, nil
		}

		defaultGroup := &types.CollectionGroup{
			Name:         "default",
			DefaultGroup: constants.CGDefault,
			UserID:       userID,
		}
		_, err = session.Insert(defaultGroup)
		if err != nil {
			return nil, err
		}
		collectionGroup = defaultGroup
		return nil, nil
	})
	if err != nil {
		return nil, errors.New(constants.DatabaseError)
	}
	return collectionGroup, nil
}

// UpdateCollectionGroup update collection group
func (cr *CollectionGroupDao) UpdateCollectionGroup(ctx context.Context, collectionGroup *types.CollectionGroup, cols []string) (err error) {
	_, err = cr.DB.Context(ctx).ID(collectionGroup.ID).Cols(cols...).Update(collectionGroup)
	if err != nil {
		return errors.New(constants.DatabaseError)
	}
	return
}

// GetCollectionGroup get collection group one
func (cr *CollectionGroupDao) GetCollectionGroup(ctx context.Context, id string) (
	collectionGroup *types.CollectionGroup, exist bool, err error,
) {
	collectionGroup = &types.CollectionGroup{}
	exist, err = cr.DB.Context(ctx).ID(id).Get(collectionGroup)
	if err != nil {
		return nil, false, errors.New(constants.DatabaseError)
	}
	return
}

// GetCollectionGroupPage get collection group page
func (cr *CollectionGroupDao) GetCollectionGroupPage(ctx context.Context, page, pageSize int, collectionGroup *types.CollectionGroup) (collectionGroupList []*types.CollectionGroup, total int64, err error) {
	collectionGroupList = make([]*types.CollectionGroup, 0)

	session := cr.DB.Context(ctx)
	if collectionGroup.UserID != "" && collectionGroup.UserID != "0" {
		session = session.Where("user_id = ?", collectionGroup.UserID)
	}
	session = session.OrderBy("update_time desc")

	total, err = tools.Help(page, pageSize, collectionGroupList, collectionGroup, session)
	err = errors.New(constants.DatabaseError)
	return
}

func (cr *CollectionGroupDao) GetDefaultID(ctx context.Context, userID string) (collectionGroup *types.CollectionGroup, has bool, err error) {
	collectionGroup = &types.CollectionGroup{}
	has, err = cr.DB.Context(ctx).Where("user_id =? and  default_group = ?", userID, constants.CGDefault).Get(collectionGroup)
	if err != nil {
		err = errors.New(constants.DatabaseError)
		return
	}
	return
}
