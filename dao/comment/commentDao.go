package comment

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/tools"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

// CommentDao comment repository
type CommentDao struct {
	DB    *xorm.Engine
	Cache *redis.Client
}

func NewCommentDao() *CommentDao {
	return &CommentDao{
		DB:    handler.Engine,
		Cache: handler.RedisClient,
	}
}

// AddComment add comment
func (cr *CommentDao) AddComment(ctx context.Context, comment *types.Comment) (err error) {
	comment.ID, err = tools.GenUniqueIDStr(ctx, comment.TableName())
	if err != nil {
		return err
	}
	_, err = cr.DB.Context(ctx).Insert(comment)
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

// RemoveComment delete comment
func (cr *CommentDao) RemoveComment(ctx context.Context, commentID string) (err error) {
	session := cr.DB.Context(ctx).ID(commentID)
	_, err = session.Update(&types.Comment{Status: constants.CommentStatusDeleted})
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

// UpdateCommentContent update comment
func (cr *CommentDao) UpdateCommentContent(
	ctx context.Context, commentID string, originalText string, parsedText string) (err error) {
	_, err = cr.DB.Context(ctx).ID(commentID).Update(&types.Comment{
		OriginalText: originalText,
		ParsedText:   parsedText,
	})
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

func (cr *CommentDao) GetCommentById(ctx context.Context, commentID string) (
	comment *types.Comment, exist bool, err error) {
	comment = &types.Comment{}
	exist, err = cr.DB.Context(ctx).ID(commentID).Get(comment)
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

func (cr *CommentDao) GetCommentCountByQuestionId(ctx context.Context, qid string) (count int64, err error) {
	list := make([]*types.Comment, 0)
	count, err = cr.DB.Context(ctx).Where("question_id = ? and status = ? ", qid, constants.CommentStatusAvailable).FindAndCount(&list)
	if err != nil {
		return count, errors.New(constants.DatabaseError)
	}
	return
}

func (cr *CommentDao) GetCommentPage(ctx context.Context, commentQuery *types.CommentQuery) (
	commentList []*types.Comment, total int64, err error) {

	commentList = make([]*types.Comment, 0)
	session := cr.DB.Context(ctx).Where("status = ?", constants.CommentStatusAvailable)
	if commentQuery.QueryCond == "vote" {
		session.OrderBy("vote_count DESC,created_at DESC")
	} else if commentQuery.QueryCond == "created_at" {
		session.OrderBy("created_at DESC")
	} else {
		session.OrderBy("created_at ASC")
	}
	cond := &types.Comment{ObjectID: commentQuery.ObjectID, UserID: commentQuery.UserID}
	total, err = tools.Help(commentQuery.Page, commentQuery.PageSize, &commentList, cond, session)
	if err != nil {
		err = errors.New(constants.DatabaseError)
	}
	return
}

func (cr *CommentDao) RemoveAllUserComment(ctx context.Context, userID string) (err error) {
	session := cr.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", constants.CommentStatusDeleted)
	affected, err := session.Update(&types.Comment{Status: constants.CommentStatusDeleted})
	if err != nil {
		return errors.New(constants.DatabaseError)
	}
	glog.Slog.Infof("delete user comment, userID: %s, affected: %d", userID, affected)
	return
}
