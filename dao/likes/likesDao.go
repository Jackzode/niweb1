package likes

import (
	"context"
	"github.com/Jackzode/painting/commons/handler"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/tools"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

type LikeDao struct {
	db    *xorm.Engine
	cache *redis.Client
}

func NewLikesRepo() *LikeDao {
	return &LikeDao{
		db:    handler.Engine,
		cache: handler.RedisClient,
	}
}

func (d *LikeDao) AddLikeRecord(ctx context.Context, likes *types.Likes) (err error) {
	likes.ID, err = tools.GenUniqueIDStr(ctx, likes.TableName())
	if err != nil {
		return err
	}
	_, err = d.db.Context(ctx).Insert(likes)
	return err
}

func (d *LikeDao) CancelLike(ctx context.Context, likes *types.Likes) error {
	//_, err := d.db.Context(ctx).Where("question_id=? and user_id=?", likes.QuestionID, likes.UserID).Cols("status", "updated_at").Update(likes)
	_, err := d.db.Context(ctx).Where("question_id=? and user_id=? and status=1", likes.QuestionID, likes.UserID).Delete(likes)
	return err
}

func (d *LikeDao) CountLikesByQuestionID(ctx context.Context, likes *types.Likes) (int64, error) {
	count, err := d.db.Context(ctx).Where("question_id=? and status=1", likes.QuestionID).Count(&types.Likes{})
	return count, err
}

func (d *LikeDao) GetLikesRecordByQuestionID(ctx context.Context, likes *types.Likes) (records []types.Likes, err error) {
	records = []types.Likes{}
	err = d.db.Context(ctx).Where("question_id=? and user_id=? and status=1", likes.QuestionID, likes.UserID).Find(&records)
	return records, err
}
