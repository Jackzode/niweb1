package collection

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/dao/tools"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

type CollectionDao struct {
	db    *xorm.Engine
	cache *redis.Client
}

func NewCollectionDao() *CollectionDao {
	return &CollectionDao{
		db:    handler.Engine,
		cache: handler.RedisClient,
	}
}

// AddCollection add collection
func (cr *CollectionDao) AddCollection(ctx context.Context, collection *types.Collection) (err error) {
	collection.ID, err = tools.GenUniqueIDStr(ctx, collection.TableName())
	if err != nil {
		return errors.New(constants.DatabaseError)
	}

	_, err = cr.db.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		old := &types.Collection{
			UserID:   collection.UserID,
			ObjectID: collection.ObjectID,
		}
		exist, err := session.ForUpdate().Get(old)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, nil
		}
		_, err = session.Insert(collection)
		if err != nil {
			return nil, err
		}
		return
	})
	if err != nil {
		return errors.New(constants.DatabaseError)
	}
	return nil
}

// RemoveCollection delete collection
func (cr *CollectionDao) RemoveCollection(ctx context.Context, id string) (err error) {
	_, err = cr.db.Context(ctx).Where("id = ?", id).Delete(&types.Collection{})
	if err != nil {
		return errors.New(constants.DatabaseError)
	}
	return nil
}

// UpdateCollection update collection
func (cr *CollectionDao) UpdateCollection(ctx context.Context, collection *types.Collection, cols []string) (err error) {
	_, err = cr.db.Context(ctx).ID(collection.ID).Cols(cols...).Update(collection)
	return errors.New(constants.DatabaseError)
}

func (cr *CollectionDao) GetCollectionById(ctx context.Context, id int) (collection *types.Collection, exist bool, err error) {
	collection = &types.Collection{}
	exist, err = cr.db.Context(ctx).ID(id).Get(collection)
	if err != nil {
		return nil, false, errors.New(constants.DatabaseError)
	}
	return
}

// GetCollectionList get collection list all
func (cr *CollectionDao) GetCollectionList(ctx context.Context, collection *types.Collection) (collectionList []*types.Collection, err error) {
	collectionList = make([]*types.Collection, 0)
	err = cr.db.Context(ctx).Find(collectionList, collection)
	err = errors.New(constants.DatabaseError)
	return
}

func (cr *CollectionDao) GetOneByObjectIDAndUser(ctx context.Context, userID string, objectID string) (collection *types.Collection, exist bool, err error) {
	collection = &types.Collection{}
	exist, err = cr.db.Context(ctx).Where("user_id = ? and object_id = ?", userID, objectID).Get(collection)
	if err != nil {
		return nil, false, errors.New(constants.DatabaseError)
	}
	return
}

// SearchByObjectIDsAndUser search by object IDs and user
func (cr *CollectionDao) SearchByObjectIDsAndUser(ctx context.Context, userID string, objectIDs []string) ([]*types.Collection, error) {
	collectionList := make([]*types.Collection, 0)
	err := cr.db.Context(ctx).Where("user_id = ?", userID).In("object_id", objectIDs).Find(&collectionList)

	if err != nil {
		return collectionList, err
	}
	return collectionList, nil
}

// CountByObjectID count by object TagID
func (cr *CollectionDao) CountByObjectID(ctx context.Context, objectID string) (total int64, err error) {
	collection := &types.Collection{}
	total, err = cr.db.Context(ctx).Where("object_id = ?", objectID).Count(collection)
	if err != nil {
		return 0, errors.New(constants.DatabaseError)
	}
	return
}

// GetCollectionPage get collection page
func (cr *CollectionDao) GetCollectionPage(ctx context.Context, page, pageSize int, collection *types.Collection) (collectionList []*types.Collection, total int64, err error) {
	collectionList = make([]*types.Collection, 0)

	session := cr.db.Context(ctx)
	if collection.UserID != "" && collection.UserID != "0" {
		session = session.Where("user_id = ?", collection.UserID)
	}

	if collection.UserCollectionGroupID != "" && collection.UserCollectionGroupID != "0" {
		session = session.Where("user_collection_group_id = ?", collection.UserCollectionGroupID)
	}
	session = session.OrderBy("update_time desc")

	total, err = tools.Help(page, pageSize, collectionList, collection, session)
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

func (cr *CollectionDao) SearchObjectCollected(ctx context.Context, userID string, objectIds []string) (map[string]bool, error) {
	for i := 0; i < len(objectIds); i++ {
		objectIds[i] = utils.DeShortID(objectIds[i])
	}

	list, err := cr.SearchByObjectIDsAndUser(ctx, userID, objectIds)
	if err != nil {
		return nil, errors.New(constants.DatabaseError)
	}

	collectedMap := make(map[string]bool)
	for _, item := range list {
		item.ObjectID = utils.EnShortID(item.ObjectID)
		collectedMap[item.ObjectID] = true
	}
	return collectedMap, nil
}

func (cr *CollectionDao) SearchList(ctx context.Context, search *types.CollectionSearch) ([]*types.Collection, int64, error) {
	var count int64
	var err error
	rows := make([]*types.Collection, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constants.DefaultPageSize
	}
	offset := search.Page * search.PageSize
	session := cr.db.Context(ctx).Where("")
	if len(search.UserID) > 0 {
		session = session.And("user_id = ?", search.UserID)
	} else {
		return rows, count, nil
	}
	session = session.Limit(search.PageSize, offset)
	count, err = session.OrderBy("updated_at desc").FindAndCount(&rows)
	if err != nil {
		return rows, count, err
	}
	return rows, count, nil
}
