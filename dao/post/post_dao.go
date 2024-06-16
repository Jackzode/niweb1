package post

import (
	"context"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/dao/tools"
	"github.com/redis/go-redis/v9"
	"time"
	"xorm.io/xorm"
)

// PostDao question repository
type PostDao struct {
	DB    *xorm.Engine
	Cache *redis.Client
}

// NewQuestionRepo new repository
func NewPostRepo() *PostDao {
	return &PostDao{
		DB:    handler.Engine,
		Cache: handler.RedisClient,
	}
}

// AddQuestion add question
func (qr *PostDao) AddQuestion(ctx context.Context, question *types.Question) (err error) {
	question.ID, err = tools.GenUniqueIDStr(ctx, question.TableName())
	if err != nil {
		return err
	}
	_, err = qr.DB.Context(ctx).Insert(question)
	if err != nil {
		return err
	}
	//todo
	question.ID = utils.EnShortID(question.ID)

	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return
}

// RemoveQuestion delete question
func (qr *PostDao) RemoveQuestion(ctx context.Context, id string) (err error) {
	id = utils.DeShortID(id)
	_, err = qr.DB.Context(ctx).Where("id =?", id).Delete(&types.Question{})
	return
}

// UpdateQuestion update question
func (qr *PostDao) UpdateQuestion(ctx context.Context, question *types.Question, Cols []string) (err error) {
	question.ID = utils.DeShortID(question.ID)
	_, err = qr.DB.Context(ctx).Where("id =?", question.ID).Cols(Cols...).Update(question)
	if err != nil {
		return err
	}
	question.ID = utils.EnShortID(question.ID)
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return
}

func (qr *PostDao) UpdatePvCount(ctx context.Context, questionID string) (err error) {
	questionID = utils.DeShortID(questionID)
	question := &types.Question{}
	_, err = qr.DB.Context(ctx).Where("id =?", questionID).Incr("view_count", 1).Update(question)
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return nil
}

func (qr *PostDao) UpdateAnswerCount(ctx context.Context, questionID string, num int) (err error) {
	questionID = utils.DeShortID(questionID)
	question := &types.Question{}
	question.AnswerCount = num
	_, err = qr.DB.Context(ctx).Where("id =?", questionID).Cols("answer_count").Update(question)
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return nil
}

func (qr *PostDao) UpdateCollectionCount(ctx context.Context, questionID string) (count int64, err error) {
	questionID = utils.DeShortID(questionID)
	_, err = qr.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		//todo &types.Collection{ObjectID: questionID}
		count, err = session.Count(1)
		if err != nil {
			return nil, err
		}

		question := &types.Question{CollectionCount: int(count)}
		_, err = session.ID(questionID).MustCols("collection_count").Update(question)
		if err != nil {
			return nil, err
		}
		return
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (qr *PostDao) UpdateQuestionStatus(ctx context.Context, questionID string, status int) (err error) {
	questionID = utils.DeShortID(questionID)
	_, err = qr.DB.Context(ctx).ID(questionID).Cols("status").Update(&types.Question{Status: status})
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, questionID)
	return nil
}

func (qr *PostDao) UpdateQuestionStatusWithOutUpdateTime(ctx context.Context, question *types.Question) (err error) {
	question.ID = utils.DeShortID(question.ID)
	_, err = qr.DB.Context(ctx).Where("id =?", question.ID).Cols("status").Update(question)
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return nil
}

func (qr *PostDao) RecoverQuestion(ctx context.Context, questionID string) (err error) {
	questionID = utils.DeShortID(questionID)
	_, err = qr.DB.Context(ctx).ID(questionID).Cols("status").Update(&types.Question{Status: types.QuestionStatusAvailable})
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, questionID)
	return nil
}

func (qr *PostDao) UpdateQuestionOperation(ctx context.Context, question *types.Question) (err error) {
	question.ID = utils.DeShortID(question.ID)
	_, err = qr.DB.Context(ctx).Where("id =?", question.ID).Cols("pin", "show").Update(question)
	return
}

func (qr *PostDao) UpdateAccepted(ctx context.Context, question *types.Question) (err error) {
	question.ID = utils.DeShortID(question.ID)
	_, err = qr.DB.Context(ctx).Where("id =?", question.ID).Cols("accepted_answer_id").Update(question)
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return nil
}

func (qr *PostDao) UpdateLastAnswer(ctx context.Context, question *types.Question) (err error) {
	question.ID = utils.DeShortID(question.ID)
	_, err = qr.DB.Context(ctx).Where("id =?", question.ID).Cols("last_answer_id").Update(question)
	if err != nil {
		return err
	}
	//todo
	//_ = qr.updateSearch(ctx, question.ID)
	return nil
}

// GetQuestion get question one
func (qr *PostDao) GetQuestion(ctx context.Context, id string) (
	question *types.Question, exist bool, err error,
) {
	id = utils.DeShortID(id)
	question = &types.Question{}
	question.ID = id
	exist, err = qr.DB.Context(ctx).Where("id = ?", id).Get(question)
	if err != nil {
		return nil, false, err
	}
	question.ID = utils.EnShortID(question.ID)
	return
}

// GetQuestionsByTitle get question list by title
func (qr *PostDao) GetQuestionsByTitle(ctx context.Context, title string, pageSize int) (
	questionList []*types.Question, err error) {
	questionList = make([]*types.Question, 0)
	session := qr.DB.Context(ctx)
	session.Where("status != ?", types.QuestionStatusDeleted)
	session.Where("title like ?", "%"+title+"%")
	session.Limit(pageSize)
	err = session.Find(&questionList)
	if err != nil {
		return nil, err
	}
	for _, item := range questionList {
		item.ID = utils.EnShortID(item.ID)
	}
	return
}

func (qr *PostDao) FindByID(ctx context.Context, id []string) (questionList []*types.Question, err error) {
	for key, itemID := range id {
		id[key] = utils.DeShortID(itemID)
	}
	questionList = make([]*types.Question, 0)
	err = qr.DB.Context(ctx).Table("question").In("id", id).Find(&questionList)
	if err != nil {
		return nil, err
	}
	for _, item := range questionList {
		item.ID = utils.EnShortID(item.ID)
	}
	return
}

// GetQuestionList get question list all
func (qr *PostDao) GetQuestionList(ctx context.Context, question *types.Question) (questionList []*types.Question, err error) {
	question.ID = utils.DeShortID(question.ID)
	questionList = make([]*types.Question, 0)
	err = qr.DB.Context(ctx).Find(questionList, question)
	if err != nil {
		return questionList, err
	}
	for _, item := range questionList {
		item.ID = utils.DeShortID(item.ID)
	}
	return
}

func (qr *PostDao) GetQuestionCount(ctx context.Context) (count int64, err error) {
	session := qr.DB.Context(ctx)
	session.In("status", []int{types.QuestionStatusAvailable, types.QuestionStatusClosed})
	count, err = session.Count(&types.Question{Show: types.QuestionShow})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (qr *PostDao) GetUserQuestionCount(ctx context.Context, userID string) (count int64, err error) {
	session := qr.DB.Context(ctx)
	session.In("status", []int{types.QuestionStatusAvailable, types.QuestionStatusClosed})
	count, err = session.Count(&types.Question{UserID: userID})
	return
}

/*  todo
func (qr *PostDao) SitemapQuestions(ctx context.Context, page, pageSize int) (
	questionIDList []*schema.SiteMapQuestionInfo, err error) {
	page = page - 1
	questionIDList = make([]*schema.SiteMapQuestionInfo, 0)

	// try to get sitemap data from cache
	cacheKey := fmt.Sprintf(constant.SiteMapQuestionCacheKeyPrefix, page)
	cacheData := qr.Cache.Get(ctx, cacheKey).String()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &questionIDList)
		return questionIDList, nil
	}

	// get sitemap data from db
	rows := make([]*types.Question, 0)
	session := qr.DB.Context(ctx)
	session.Select("id,title,created_at,post_update_time")
	session.Where("`show` = ?", types.QuestionShow)
	session.Where("status = ? OR status = ?", types.QuestionStatusAvailable, types.QuestionStatusClosed)
	session.Limit(pageSize, page*pageSize)
	session.Asc("created_at")
	err = session.Find(&rows)
	if err != nil {
		return questionIDList, err
	}

	// warp data
	for _, question := range rows {
		item := &schema.SiteMapQuestionInfo{ID: question.ID}
		if utils.GetEnableShortID(ctx) {
			item.ID = utils.EnShortID(question.ID)
		}
		item.Title = htmltext.UrlTitle(question.Title)
		if question.PostUpdateTime.IsZero() {
			item.UpdateTime = question.CreatedAt.Format(time.RFC3339)
		} else {
			item.UpdateTime = question.PostUpdateTime.Format(time.RFC3339)
		}
		questionIDList = append(questionIDList, item)
	}

	// set sitemap data to cache
	cacheDataByte, _ := json.Marshal(questionIDList)
	if err := qr.Cache.Set(ctx, cacheKey, string(cacheDataByte), constant.SiteMapQuestionCacheTime).Err(); err != nil {
		glog.Slog.Error(err)
	}
	return questionIDList, nil
}
*/

// GetQuestionPage query question page
func (qr *PostDao) GetQuestionPage(ctx context.Context, page, pageSize int, userID, tagID, orderCond string, inDays int) (
	questionList []*types.Question, total int64, err error) {
	questionList = make([]*types.Question, 0)

	session := qr.DB.Context(ctx).Where("question.status = ? OR question.status = ?",
		types.QuestionStatusAvailable, types.QuestionStatusClosed)
	if len(tagID) > 0 {
		session.Join("LEFT", "tag_rel", "question.id = tag_rel.object_id")
		session.And("tag_rel.tag_id = ?", tagID)
		//todo types.TagRelStatusAvailable
		session.And("tag_rel.status = ?", 1)
	}
	if len(userID) > 0 {
		session.And("question.user_id = ?", userID)
	} else {
		session.And("question.show = ?", types.QuestionShow)
	}
	if inDays > 0 {
		session.And("question.created_at > ?", time.Now().AddDate(0, 0, -inDays))
	}

	switch orderCond {
	case "newest":
		session.OrderBy("question.pin desc,question.created_at DESC")
	case "active":
		session.OrderBy("question.pin desc,question.post_update_time DESC, question.updated_at DESC")
	case "frequent":
		session.OrderBy("question.pin desc,question.view_count DESC")
	case "score":
		session.OrderBy("question.pin desc,question.vote_count DESC, question.view_count DESC")
	case "unanswered":
		session.Where("question.last_answer_id = 0")
		session.OrderBy("question.pin desc,question.created_at DESC")
	}

	total, err = tools.Help(page, pageSize, &questionList, &types.Question{}, session)
	for _, item := range questionList {
		item.ID = utils.EnShortID(item.ID)
	}
	return questionList, total, err
}

/*
func (qr *PostDao) AdminQuestionPage(ctx context.Context, search *schema.AdminQuestionPageReq) ([]*types.Question, int64, error) {
	var (
		count   int64
		err     error
		session = qr.DB.Context(ctx).Table("question")
	)

	session.Where(builder.Eq{
		"status": search.Status,
	})

	rows := make([]*types.Question, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}

	// search by question title like or question id
	if len(search.Query) > 0 {
		// check id search
		var (
			idSearch = false
			id       = ""
		)

		if strings.Contains(search.Query, "question:") {
			idSearch = true
			id = strings.TrimSpace(strings.TrimPrefix(search.Query, "question:"))
			id = utils.DeShortID(id)
			for _, r := range id {
				if !unicode.IsDigit(r) {
					idSearch = false
					break
				}
			}
		}

		if idSearch {
			session.And(builder.Eq{
				"id": id,
			})
		} else {
			session.And(builder.Like{
				"title", search.Query,
			})
		}
	}

	offset := search.Page * search.PageSize

	session.OrderBy("created_at desc").
		Limit(search.PageSize, offset)
	count, err = session.FindAndCount(&rows)
	if err != nil {
		return rows, count, err
	}
	if utils.GetEnableShortID(ctx) {
		for _, item := range rows {
			item.ID = utils.EnShortID(item.ID)
		}
	}
	return rows, count, nil
}



// updateSearch update search, if search plugin not enable, do nothing
func (qr *PostDao) updateSearch(ctx context.Context, questionID string) (err error) {
	// check search plugin
	var s plugin.Search
	_ = plugin.CallSearch(func(search plugin.Search) error {
		s = search
		return nil
	})
	if s == nil {
		return
	}
	questionID = utils.DeShortID(questionID)
	question, exist, err := qr.GetQuestion(ctx, questionID)
	if !exist {
		return
	}
	if err != nil {
		return err
	}

	// get tags
	var (
		tagListList = make([]*types.TagRel, 0)
		tags        = make([]string, 0)
	)
	session := qr.DB.Context(ctx).Where("object_id = ?", questionID)
	session.Where("status = ?", types.TagRelStatusAvailable)
	err = session.Find(&tagListList)
	if err != nil {
		return
	}
	for _, tag := range tagListList {
		tags = append(tags, tag.TagID)
	}
	content := &plugin.SearchContent{
		ObjectID:    questionID,
		Title:       question.Title,
		Type:        constant.QuestionObjectType,
		Content:     question.OriginalText,
		Answers:     int64(question.AnswerCount),
		Status:      plugin.SearchContentStatus(question.Status),
		Tags:        tags,
		QuestionID:  questionID,
		UserID:      question.UserID,
		Views:       int64(question.ViewCount),
		Created:     question.CreatedAt.Unix(),
		Active:      question.UpdatedAt.Unix(),
		Score:       int64(question.VoteCount),
		HasAccepted: question.AcceptedAnswerID != "" && question.AcceptedAnswerID != "0",
	}
	err = s.UpdateContent(ctx, content)
	return
}

*/

func (qr *PostDao) RemoveAllUserQuestion(ctx context.Context, userID string) (err error) {
	// get all question id that need to be deleted
	questionIDs := make([]string, 0)
	session := qr.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", types.QuestionStatusDeleted)
	err = session.Select("id").Table("question").Find(&questionIDs)
	if err != nil {
		return err
	}
	if len(questionIDs) == 0 {
		return nil
	}

	glog.Slog.Infof("find %d questions need to be deleted for user %s", len(questionIDs), userID)

	// delete all question
	session = qr.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", types.QuestionStatusDeleted)
	_, err = session.Cols("status", "updated_at").Update(&types.Question{
		UpdatedAt: time.Now(),
		Status:    types.QuestionStatusDeleted,
	})
	if err != nil {
		return err
	}

	//todo update search content
	//for _, id := range questionIDs {

	//_ = qr.updateSearch(ctx, id)
	//}
	return nil
}
