package likes

import (
	"context"
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/likes"
	"github.com/gin-gonic/gin"
	"time"
)

type LikeService struct {
	likeDao *likes.LikeDao
}

func NewLikeService() *LikeService {

	return &LikeService{
		likeDao: likes.NewLikesRepo(),
	}
}

func (s *LikeService) AddLike(ctx context.Context, req *types.AddLikeReq) error {
	bean := &types.Likes{}
	bean.UserID = req.UserId
	bean.AuthorID = req.AuthorID
	bean.QuestionID = req.QuestionID
	bean.Status = constants.Like
	bean.CreatedAt = time.Now()
	err := s.likeDao.AddLikeRecord(ctx, bean)
	if err != nil {
		glog.Slog.Error(err)
	}
	return err
}

func (s *LikeService) CancelLike(ctx context.Context, req *types.AddLikeReq) error {
	bean := &types.Likes{}
	bean.UserID = req.UserId
	bean.AuthorID = req.AuthorID
	bean.QuestionID = req.QuestionID
	//bean.Status = constants.Dislike
	//bean.UpdatedAt = time.Now()
	err := s.likeDao.CancelLike(ctx, bean)
	if err != nil {
		glog.Slog.Error(err)
	}
	return err
}

func (s *LikeService) CountLikes(ctx context.Context, req *types.AddLikeReq) (int64, error) {
	bean := &types.Likes{}
	bean.QuestionID = req.QuestionID
	bean.UserID = req.UserId
	count, err := s.likeDao.CountLikesByQuestionID(ctx, bean)
	if err != nil {
		glog.Slog.Error(err)
		return 0, err
	}
	return count, nil
}

func (s *LikeService) CheckLiked(ctx *gin.Context, req *types.AddLikeReq) (bool, error) {
	bean := &types.Likes{}
	bean.QuestionID = req.QuestionID
	bean.UserID = req.UserId
	records, err := s.likeDao.GetLikesRecordByQuestionID(ctx, bean)
	if err != nil {
		glog.Slog.Error(err)
		return false, err
	}
	if len(records) != 1 {
		return false, nil
	}
	if records[0].UserID != req.UserId {
		return false, nil
	}
	return true, nil
}
