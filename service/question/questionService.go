package question

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/dao/post"
	collectService "github.com/Jackzode/painting/service/collection"
	"github.com/Jackzode/painting/service/user"
	"time"
)

type QuestionService struct {
	questionDao    *post.PostDao
	userService    *user.UserService
	collectService *collectService.CollectionService
}

func NewQuestionService() *QuestionService {

	return &QuestionService{
		questionDao:    post.NewPostRepo(),
		userService:    user.NewUserService(),
		collectService: collectService.NewCollectionService(),
	}
}

// PersonalCollectionPage get collection list by user
func (qs *QuestionService) PersonalCollectionPage(ctx context.Context, req *types.PersonalCollectionPageReq) (
	pageModel *utils.PageModel, err error) {
	list := make([]*types.QuestionPageResp, 0)
	collectionSearch := &types.CollectionSearch{}
	collectionSearch.UserID = req.UserID
	collectionSearch.Page = req.Page
	collectionSearch.PageSize = req.PageSize
	collectionList, total, err := qs.collectService.SearchList(ctx, collectionSearch)
	if err != nil {
		return nil, err
	}
	questionIDs := make([]string, 0)
	for _, item := range collectionList {
		questionIDs = append(questionIDs, item.ObjectID)
	}

	questionMaps, err := qs.FindInfoByIDs(ctx, questionIDs, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, id := range questionIDs {
		_, ok := questionMaps[id]
		if ok {
			//questionMaps[id].LastAnsweredUserInfo = nil
			//questionMaps[id].UpdateUserInfo = nil
			//questionMaps[id].Content = ""
			//questionMaps[id].HTML = ""
			list = append(list, questionMaps[id])
		}
	}

	return utils.NewPageModel(total, list), nil
}

// PersonalQuestionPage get question list by user
func (qs *QuestionService) PersonalQuestionPage(ctx context.Context, req *types.PersonalQuestionPageReq) (
	pageModel *utils.PageModel, err error) {

	//userinfo, exist, err := qs.userService.GetUserBasicInfoByUserName(ctx, req.Username)
	//if err != nil {
	//	return nil, err
	//}
	//if !exist {
	//	return nil, errors.New(constants.UserNotFound)
	//}
	search := &types.QuestionPageReq{}
	search.OrderCond = req.OrderCond
	search.Page = req.Page
	search.PageSize = req.PageSize
	search.Username = req.Username
	search.UserIDBeSearched = req.LoginUserID
	search.LoginUserID = req.LoginUserID
	questionList, total, err := qs.GetQuestionPage(ctx, search)
	if err != nil {
		return nil, err
	}
	return utils.NewPageModel(total, questionList), nil
}

// UpdateQuestion update question
func (qs *QuestionService) UpdateQuestion(ctx context.Context, req *types.QuestionUpdate) (questionInfo any, err error) {

	questionInfo = &types.QuestionInfo{}
	dbinfo, has, err := qs.questionDao.GetQuestion(ctx, req.ID)
	if err != nil {
		glog.Slog.Error(err.Error())
		return
	}
	if !has {
		err = errors.New("question not found")
		return
	}
	if dbinfo.Status == constants.QuestionStatusDeleted {
		err = errors.New(constants.QuestionCannotUpdate)
		return nil, err
	}

	now := time.Now()
	question := &types.Question{}
	question.Title = req.Title
	question.OriginalText = req.Content
	question.ParsedText = req.HTML
	question.ID = utils.DeShortID(req.ID)
	question.UpdatedAt = now
	question.PostUpdateTime = now
	question.UserID = dbinfo.UserID
	question.LastEditUserID = req.UserID

	//If the content is the same, ignore it
	if dbinfo.Title == req.Title && dbinfo.OriginalText == req.Content {
		return
	}

	questionInfo, err = qs.GetQuestion(ctx, question.ID, question.UserID)
	return
}

// GetQuestionPage query questions page
func (qs *QuestionService) GetQuestionPage(ctx context.Context, req *types.QuestionPageReq) (
	questions []*types.QuestionPageResp, total int64, err error) {

	questions = make([]*types.QuestionPageResp, 0)
	// query by user condition
	if req.Username != "" {
		userinfo, exist, err := qs.userService.GetUserBasicInfoByUserName(ctx, req.Username)
		if err != nil {
			return nil, 0, err
		}
		if !exist {
			return questions, 0, nil
		}
		req.UserIDBeSearched = userinfo.ID
	}

	questionList, total, err := qs.questionDao.GetQuestionPage(ctx, req.Page, req.PageSize,
		req.UserIDBeSearched, req.TagID, req.OrderCond, req.InDays)
	if err != nil {
		return nil, 0, err
	}
	questions, err = qs.FormatQuestionsPage(ctx, questionList, req.LoginUserID, req.OrderCond)
	if err != nil {
		return nil, 0, err
	}
	return questions, total, nil
}

// AddQuestion add question
func (qs *QuestionService) AddQuestion(ctx context.Context, req *types.QuestionAdd) (questionInfo any, err error) {

	question := &types.Question{}
	now := time.Now()
	question.UserID = req.UserID
	question.Title = req.Title
	question.OriginalText = req.Content
	question.ParsedText = req.HTML
	question.AcceptedAnswerID = "0"
	question.LastAnswerID = "0"
	question.LastEditUserID = "0"
	question.AllowReprint = utils.Str2Int(req.AllowReprint)
	question.AllowComment = utils.Str2Int(req.AllowComment)
	question.CopyRight = utils.Str2Int(req.CopyRight)
	question.Feeds = utils.Str2Int(req.Feeds)
	question.Status = constants.QuestionStatusAvailable
	question.RevisionID = "0"
	question.CreatedAt = now
	question.PostUpdateTime = now
	question.Pin = constants.QuestionUnPin
	question.Show = constants.QuestionShow
	//question.UpdatedAt = nil
	err = qs.questionDao.AddQuestion(ctx, question)
	if err != nil {
		glog.Slog.Error(err.Error())
		return
	}

	// user add question count
	userQuestionCount, err := qs.GetUserQuestionCount(ctx, question.UserID)
	if err != nil {
		glog.Slog.Errorf("get user question count error %v", err)
	} else {
		err = qs.userService.UpdateQuestionCount(ctx, question.UserID, userQuestionCount)
		if err != nil {
			glog.Slog.Errorf("update user question count error %v", err)
		}
	}
	questionInfo, err = qs.GetQuestion(ctx, question.ID, question.UserID)
	return
}

func (qs *QuestionService) GetQuestionAndAddPV(ctx context.Context, questionID, loginUserID string) (resp *types.QuestionInfo, err error) {

	err = qs.UpdatePv(ctx, questionID)
	if err != nil {
		glog.Slog.Error(err)
	}
	return qs.GetQuestion(ctx, questionID, loginUserID)
}

// GetQuestion get question one
func (qs *QuestionService) GetQuestion(ctx context.Context, questionID, userID string) (resp *types.QuestionInfo, err error) {

	question, err := qs.Info(ctx, questionID, userID)
	if err != nil {
		return
	}
	//将html中的标签去掉，只保留240个字符展示，多余的用。。。代替
	question.Description = utils.FetchExcerpt(question.HTML, "...", 240)
	return question, nil
}

func (qs *QuestionService) GetUserQuestionCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.questionDao.GetUserQuestionCount(ctx, userID)
}

func (qs *QuestionService) UpdatePv(ctx context.Context, questionID string) error {
	return qs.questionDao.UpdatePvCount(ctx, questionID)
}

func (qs *QuestionService) UpdateCollectionCount(ctx context.Context, questionID string) (count int64, err error) {
	return qs.questionDao.UpdateCollectionCount(ctx, questionID)
}

func (qs *QuestionService) UpdateAccepted(ctx context.Context, questionID, AnswerID string) error {
	question := &types.Question{}
	question.ID = questionID
	question.AcceptedAnswerID = AnswerID
	return qs.questionDao.UpdateAccepted(ctx, question)
}

func (qs *QuestionService) UpdateLastAnswer(ctx context.Context, questionID, AnswerID string) error {
	question := &types.Question{}
	question.ID = questionID
	question.LastAnswerID = AnswerID
	return qs.questionDao.UpdateLastAnswer(ctx, question)
}

func (qs *QuestionService) UpdatePostTime(ctx context.Context, questionID string) error {
	questioninfo := &types.Question{}
	now := time.Now()
	questioninfo.ID = questionID
	questioninfo.PostUpdateTime = now
	return qs.questionDao.UpdateQuestion(ctx, questioninfo, []string{"post_update_time"})
}

func (qs *QuestionService) UpdatePostSetTime(ctx context.Context, questionID string, setTime time.Time) error {
	questioninfo := &types.Question{}
	questioninfo.ID = questionID
	questioninfo.PostUpdateTime = setTime
	return qs.questionDao.UpdateQuestion(ctx, questioninfo, []string{"post_update_time"})
}

func (qs *QuestionService) FindInfoByIDs(ctx context.Context, questionIDs []string, loginUserID string) (map[string]*types.QuestionPageResp, error) {
	list := make(map[string]*types.QuestionPageResp)
	questionList, err := qs.questionDao.FindByID(ctx, questionIDs)
	if err != nil {
		return list, err
	}
	questions, err := qs.FormatQuestionsPage(ctx, questionList, loginUserID, "")
	if err != nil {
		return list, err
	}
	for _, item := range questions {
		list[item.ID] = item
	}
	return list, nil
}

func (qs *QuestionService) InviteUserInfo(ctx context.Context, questionID string) (inviteList []*types.UserBasicInfo, err error) {
	InviteUserInfo := make([]*types.UserBasicInfo, 0)
	dbinfo, has, err := qs.questionDao.GetQuestion(ctx, questionID)
	if err != nil {
		return InviteUserInfo, err
	}
	if !has {
		return InviteUserInfo, errors.New(constants.QuestionNotFound)
	}
	//InviteUser
	if dbinfo.InviteUserID != "" {
		InviteUserIDs := make([]string, 0)
		err := json.Unmarshal([]byte(dbinfo.InviteUserID), &InviteUserIDs)
		if err == nil {
			//inviteUserInfoMap, err := UserCommonServicer.BatchUserBasicInfoByID(ctx, InviteUserIDs)
			//if err == nil {
			//	for _, userid := range InviteUserIDs {
			//		_, ok := inviteUserInfoMap[userid]
			//		if ok {
			//			InviteUserInfo = append(InviteUserInfo, inviteUserInfoMap[userid])
			//		}
			//	}
			//}
		}
	}
	return InviteUserInfo, nil
}

func (qs *QuestionService) Info(ctx context.Context, questionID string, loginUserID string) (showinfo *types.QuestionInfo, err error) {

	dbinfo, has, err := qs.questionDao.GetQuestion(ctx, questionID)
	if err != nil {
		return showinfo, err
	}
	dbinfo.ID = utils.DeShortID(dbinfo.ID)
	if !has {
		return showinfo, errors.New(constants.QuestionNotFound)
	}
	showinfo = qs.ShowFormat(ctx, dbinfo)

	userIds := make([]string, 0)
	if utils.IsNotZeroString(dbinfo.UserID) {
		userIds = append(userIds, dbinfo.UserID)
	}
	if utils.IsNotZeroString(dbinfo.LastEditUserID) {
		userIds = append(userIds, dbinfo.LastEditUserID)
	}
	if utils.IsNotZeroString(showinfo.LastAnsweredUserID) {
		userIds = append(userIds, showinfo.LastAnsweredUserID)
	}
	userInfoMap, err := qs.userService.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return showinfo, err
	}

	_, ok := userInfoMap[dbinfo.UserID]
	if ok {
		showinfo.UserInfo = userInfoMap[dbinfo.UserID]
	}
	_, ok = userInfoMap[dbinfo.LastEditUserID]
	if ok {
		showinfo.UpdateUserInfo = userInfoMap[dbinfo.LastEditUserID]
	}
	_, ok = userInfoMap[showinfo.LastAnsweredUserID]
	if ok {
		showinfo.LastAnsweredUserInfo = userInfoMap[showinfo.LastAnsweredUserID]
	}

	if loginUserID == "" {
		return showinfo, nil
	}

	//showinfo.VoteStatus = qs.questionDao.VoteRepo.GetVoteStatus(ctx, questionID, loginUserID)

	// // check is followed
	//isFollowed, _ := qs.questionDao.FollowRepo.IsFollowed(ctx, loginUserID, questionID)
	//showinfo.IsFollowed = isFollowed

	//collectedMap, err := CollectionCommon.SearchObjectCollected(ctx, loginUserID, []string{dbinfo.ID})
	//if err != nil {
	//	return nil, err
	//}
	//if len(collectedMap) > 0 {
	//	showinfo.Collected = true
	//}
	return showinfo, nil
}

func (qs *QuestionService) FormatQuestionsPage(
	ctx context.Context, questionList []*types.Question, loginUserID string, orderCond string) (
	formattedQuestions []*types.QuestionPageResp, err error) {

	formattedQuestions = make([]*types.QuestionPageResp, 0)
	userIdMap := make(map[string]struct{})
	userIDs := make([]string, 0)
	for _, questionInfo := range questionList {
		t := &types.QuestionPageResp{
			ID:               questionInfo.ID,
			CreatedAt:        questionInfo.CreatedAt.Unix(),
			Title:            questionInfo.Title,
			UrlTitle:         utils.UrlTitle(questionInfo.Title),
			Description:      utils.FetchExcerpt(questionInfo.ParsedText, "...", 240),
			Status:           questionInfo.Status,
			ViewCount:        questionInfo.ViewCount,
			UniqueViewCount:  questionInfo.UniqueViewCount,
			VoteCount:        questionInfo.VoteCount,
			AnswerCount:      questionInfo.AnswerCount,
			CollectionCount:  questionInfo.CollectionCount,
			FollowCount:      questionInfo.FollowCount,
			AcceptedAnswerID: questionInfo.AcceptedAnswerID,
			LastAnswerID:     questionInfo.LastAnswerID,
			Pin:              questionInfo.Pin,
			Show:             questionInfo.Show,
			AuthorID:         questionInfo.UserID,
			Content:          questionInfo.OriginalText,
		}

		//questionIDs = append(questionIDs, questionInfo.ID)
		//userIDs = append(userIDs, questionInfo.UserID)
		userIdMap[questionInfo.UserID] = struct{}{}
		formattedQuestions = append(formattedQuestions, t)
	}
	for id := range userIdMap {
		userIDs = append(userIDs, id)
	}
	basicInfoByID, err := qs.userService.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return formattedQuestions, err
	}
	for _, question := range formattedQuestions {
		useInfo, ok := basicInfoByID[question.AuthorID]
		if ok {
			question.AuthorInfo.Avatar = useInfo.Avatar
			question.AuthorInfo.Username = useInfo.Username
			question.AuthorInfo.Rank = useInfo.Rank
			question.AuthorInfo.Description = useInfo.Description
		}
	}

	return formattedQuestions, nil
}

func (qs *QuestionService) FormatQuestions(ctx context.Context, questionList []*types.Question, loginUserID string) ([]*types.QuestionInfo, error) {
	list := make([]*types.QuestionInfo, 0)
	objectIds := make([]string, 0)
	userIds := make([]string, 0)

	for _, questionInfo := range questionList {
		item := qs.ShowFormat(ctx, questionInfo)
		list = append(list, item)
		objectIds = append(objectIds, item.ID)
		userIds = append(userIds, item.UserID, item.LastEditUserID, item.LastAnsweredUserID)
	}
	//tagsMap, err := TagServicer.BatchGetObjectTag(ctx, objectIds)
	//if err != nil {
	//	return list, err
	//}
	//
	//userInfoMap, err := UserCommonServicer.BatchUserBasicInfoByID(ctx, userIds)
	//if err != nil {
	//	return list, err
	//}

	//for _, item := range list {
	//	item.Tags = tagsMap[item.ID]
	//	item.UserInfo = userInfoMap[item.UserID]
	//	item.UpdateUserInfo = userInfoMap[item.LastEditUserID]
	//	item.LastAnsweredUserInfo = userInfoMap[item.LastAnsweredUserID]
	//}
	if loginUserID == "" {
		return list, nil
	}

	//collectedMap, err := CollectionCommon.SearchObjectCollected(ctx, loginUserID, objectIds)
	//if err != nil {
	//	return nil, err
	//}
	//for _, item := range list {
	//	item.Collected = collectedMap[item.ID]
	//}
	return list, nil
}

// RemoveQuestion delete question
func (qs *QuestionService) RemoveQuestion(ctx context.Context, req *types.RemoveQuestionReq) (err error) {
	questionInfo, has, err := qs.questionDao.GetQuestion(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	if questionInfo.Status == constants.QuestionStatusDeleted {
		return nil
	}

	questionInfo.Status = constants.QuestionStatusDeleted
	err = qs.questionDao.UpdateQuestionStatus(ctx, questionInfo.ID, questionInfo.Status)
	if err != nil {
		return err
	}

	//userQuestionCount, err := qs.GetUserQuestionCount(ctx, questionInfo.UserID)
	//if err != nil {
	//	glog.Slog.Error("user GetUserQuestionCount error", err.Error())
	//} else {
	//err = UserCommonServicer.UpdateQuestionCount(ctx, questionInfo.UserID, userQuestionCount)
	//if err != nil {
	//	glog.Slog.Error("user IncreaseQuestionCount error", err.Error())
	//}
	//}

	return nil
}

func (qs *QuestionService) CloseQuestion(ctx context.Context, req *types.CloseQuestionReq) error {
	questionInfo, has, err := qs.questionDao.GetQuestion(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	questionInfo.Status = constants.QuestionStatusClosed
	err = qs.questionDao.UpdateQuestionStatus(ctx, questionInfo.ID, questionInfo.Status)
	if err != nil {
		return err
	}

	//closeMeta, _ := json.Marshal(types.CloseQuestionMeta{
	//	CloseType: req.CloseType,
	//	CloseMsg:  req.CloseMsg,
	//})
	//err = MetaService.AddMeta(ctx, req.ID, types.QuestionCloseReasonKey, string(closeMeta))
	//if err != nil {
	//	return err
	//}
	//
	//ActivityQueueServicer.Send(ctx, &types.ActivityMsg{
	//	UserID:           questionInfo.UserID,
	//	ObjectID:         questionInfo.ID,
	//	OriginalObjectID: questionInfo.ID,
	//	ActivityTypeKey:  constant.ActQuestionClosed,
	//})
	return nil
}

/*
func (as *QuestionService) RemoveAnswer(ctx context.Context, id string) (err error) {
	answerinfo, has, err := qs.questionDao.AnswerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	// user add question count

	err = as.UpdateAnswerCount(ctx, answerinfo.QuestionID)
	if err != nil {
		log.Error("UpdateAnswerCount error", err.Error())
	}
	userAnswerCount, err := qs.questionDao.AnswerRepo.GetCountByUserID(ctx, answerinfo.UserID)
	if err != nil {
		log.Error("GetCountByUserID error", err.Error())
	}
	err = UserCommonServicer.UpdateAnswerCount(ctx, answerinfo.UserID, int(userAnswerCount))
	if err != nil {
		log.Error("user UpdateAnswerCount error", err.Error())
	}

	return qs.questionDao.AnswerRepo.RemoveAnswer(ctx, id)
}

*/

func (qs *QuestionService) SitemapCron(ctx context.Context) {
	//questionNum, err := qs.questionDao.GetQuestionCount(ctx)
	//if err != nil {
	//	glog.Slog.Error(err)
	//	return
	//}
	//if questionNum <= constants.SitemapMaxSize {
	//	_, err = qs.questionDao.SitemapQuestions(ctx, 1, int(questionNum))
	//	if err != nil {
	//		glog.Slog.Errorf("get site map question error: %v", err)
	//	}
	//	return
	//}

	//totalPages := int(math.Ceil(float64(questionNum) / float64(constant.SitemapMaxSize)))
	//for i := 1; i <= totalPages; i++ {
	//	_, err = qs.questionDao.SitemapQuestions(ctx, i, constant.SitemapMaxSize)
	//	if err != nil {
	//		log.Errorf("get site map question error: %v", err)
	//		return
	//	}
	//}
}

func (qs *QuestionService) SetCache(ctx context.Context, cachekey string, info interface{}) error {
	//infoStr, err := json.Marshal(info)
	//if err != nil {
	//	return errors.New(constants.UnknownError)
	//}

	//err = handler.RedisClient.Set(ctx, cachekey, string(infoStr), types.DashboardCacheTime).Err()
	//if err != nil {
	//	return errors.New(constants.UnknownError)
	//}
	return nil
}

func (qs *QuestionService) ShowListFormat(ctx context.Context, data *types.Question) *types.QuestionInfo {
	return qs.ShowFormat(ctx, data)
}

func (qs *QuestionService) ShowFormat(ctx context.Context, data *types.Question) *types.QuestionInfo {
	info := types.QuestionInfo{}
	info.ID = data.ID
	//if utils.GetEnableShortID(ctx) {
	//	info.ID = uid.EnShortID(data.ID)
	//}
	info.Title = data.Title
	//info.UrlTitle = htmltext.UrlTitle(data.Title)
	info.Content = data.OriginalText
	info.HTML = data.ParsedText
	info.ViewCount = data.ViewCount
	info.UniqueViewCount = data.UniqueViewCount
	info.VoteCount = data.VoteCount
	info.AnswerCount = data.AnswerCount
	info.CollectionCount = data.CollectionCount
	info.FollowCount = data.FollowCount
	info.AcceptedAnswerID = data.AcceptedAnswerID
	info.LastAnswerID = data.LastAnswerID
	info.CreateTime = data.CreatedAt.Unix()
	info.UpdateTime = data.UpdatedAt.Unix()
	info.PostUpdateTime = data.PostUpdateTime.Unix()
	if data.PostUpdateTime.Unix() < 1 {
		info.PostUpdateTime = 0
	}
	info.QuestionUpdateTime = data.UpdatedAt.Unix()
	if data.UpdatedAt.Unix() < 1 {
		info.QuestionUpdateTime = 0
	}
	info.Status = data.Status
	info.Pin = data.Pin
	info.Show = data.Show
	info.UserID = data.UserID
	info.LastEditUserID = data.LastEditUserID
	//if data.LastAnswerID != "0" {
	//	answerInfo, exist, err := qs.questionDao.AnswerRepo.GetAnswer(ctx, data.LastAnswerID)
	//	if err == nil && exist {
	//		if answerInfo.LastEditUserID != "0" {
	//			info.LastAnsweredUserID = answerInfo.LastEditUserID
	//		} else {
	//			info.LastAnsweredUserID = answerInfo.UserID
	//		}
	//	}
	//
	//}
	//info.Tags = make([]*types.TagResp, 0)
	return &info
}

//func (qs *QuestionService) ShowFormatWithTag(ctx context.Context, data *types.QuestionWithTagsRevision) *types.QuestionInfo {
//	info := qs.ShowFormat(ctx, &data.Question)
//	Tags := make([]*types.TagResp, 0)
//	for _, tag := range data.Tags {
//		item := &types.TagResp{}
//		item.SlugName = tag.SlugName
//		item.DisplayName = tag.DisplayName
//		item.Recommend = tag.Recommend
//		item.Reserved = tag.Reserved
//		Tags = append(Tags, item)
//	}
//	info.Tags = Tags
//	return info
//}
