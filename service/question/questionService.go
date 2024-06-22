package question

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/post"
	"time"
)

type QuestionService struct {
	questionDao *post.PostDao
}

func NewQuestionService() *QuestionService {

	return &QuestionService{
		questionDao: post.NewPostRepo(),
	}
}

func (qs *QuestionService) GetUserQuestionCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.questionDao.GetUserQuestionCount(ctx, userID)
}

func (qs *QuestionService) UpdatePv(ctx context.Context, questionID string) error {
	return qs.questionDao.UpdatePvCount(ctx, questionID)
}

func (qs *QuestionService) UpdateAnswerCount(ctx context.Context, questionID string) error {
	//count, err := qs.questionDao.AnswerRepo.GetCountByQuestionID(ctx, questionID)
	//if err != nil {
	//	return err
	//}
	//return qs.questionDao.UpdateAnswerCount(ctx, questionID, int(count))
	return nil
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

func (qs *QuestionService) FindInfoByID(ctx context.Context, questionIDs []string, loginUserID string) (map[string]*types.QuestionInfo, error) {
	list := make(map[string]*types.QuestionInfo)
	questionList, err := qs.questionDao.FindByID(ctx, questionIDs)
	if err != nil {
		return list, err
	}
	questions, err := qs.FormatQuestions(ctx, questionList, loginUserID)
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

/*
func (qs *QuestionService) Info(ctx context.Context, questionID string, loginUserID string) (showinfo *types.QuestionInfo, err error) {

	dbinfo, has, err := qs.questionDao.GetQuestion(ctx, questionID)
	if err != nil {
		return showinfo, err
	}
	dbinfo.ID = uid.DeShortID(dbinfo.ID)
	if !has {
		return showinfo, errors.New(constants.QuestionNotFound)
	}
	showinfo = qs.ShowFormat(ctx, dbinfo)

	//
	if showinfo.Status == 2 {
		var metainfo *types.Meta
		metainfo, err = MetaService.GetMetaByObjectIdAndKey(ctx, dbinfo.ID, types.QuestionCloseReasonKey)
		if err != nil {
			glog.Slog.Error(err)
		} else {
			// metainfo.Value
			closemsg := &types.CloseQuestionMeta{}
			err = json.Unmarshal([]byte(metainfo.Value), closemsg)
			if err != nil {
				glog.Slog.Error("json.Unmarshal CloseQuestionMeta error", err.Error())
			} else {
				cfg, err := utils.GetConfigByID(ctx, closemsg.CloseType)
				if err != nil {
					glog.Slog.Error("json.Unmarshal QuestionCloseJson error", err.Error())
				} else {
					reasonItem := &types.ReasonItem{}
					_ = json.Unmarshal(cfg.GetByteValue(), reasonItem)
					reasonItem.Translate(cfg.Key, utils.GetLangByCtx(ctx))
					operation := &types.Operation{}
					operation.Type = reasonItem.Name
					operation.Description = reasonItem.Description
					operation.Msg = closemsg.CloseMsg
					operation.Time = metainfo.CreatedAt.Unix()
					operation.Level = types.OperationLevelInfo
					showinfo.Operation = operation
				}
			}
		}
	}

	tagmap, err := TagServicer.GetObjectTag(ctx, questionID)
	if err != nil {
		return showinfo, err
	}
	showinfo.Tags = tagmap

	userIds := make([]string, 0)
	if checker.IsNotZeroString(dbinfo.UserID) {
		userIds = append(userIds, dbinfo.UserID)
	}
	if checker.IsNotZeroString(dbinfo.LastEditUserID) {
		userIds = append(userIds, dbinfo.LastEditUserID)
	}
	if checker.IsNotZeroString(showinfo.LastAnsweredUserID) {
		userIds = append(userIds, showinfo.LastAnsweredUserID)
	}
	userInfoMap, err := UserCommonServicer.BatchUserBasicInfoByID(ctx, userIds)
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

	showinfo.VoteStatus = qs.questionDao.VoteRepo.GetVoteStatus(ctx, questionID, loginUserID)

	// // check is followed
	isFollowed, _ := qs.questionDao.FollowRepo.IsFollowed(ctx, loginUserID, questionID)
	showinfo.IsFollowed = isFollowed

	ids, err := AnswerCommonServicer.SearchAnswerIDs(ctx, loginUserID, dbinfo.ID)
	if err != nil {
		glog.Slog.Error("AnswerFunc.SearchAnswerIDs", err)
	}
	showinfo.Answered = len(ids) > 0
	if showinfo.Answered {
		showinfo.FirstAnswerId = ids[0]
	}

	collectedMap, err := CollectionCommon.SearchObjectCollected(ctx, loginUserID, []string{dbinfo.ID})
	if err != nil {
		return nil, err
	}
	if len(collectedMap) > 0 {
		showinfo.Collected = true
	}
	return showinfo, nil
}
*/

func (qs *QuestionService) FormatQuestionsPage(
	ctx context.Context, questionList []*types.Question, loginUserID string, orderCond string) (
	formattedQuestions []*types.QuestionPageResp, err error) {
	formattedQuestions = make([]*types.QuestionPageResp, 0)
	questionIDs := make([]string, 0)
	userIDs := make([]string, 0)
	for _, questionInfo := range questionList {
		t := &types.QuestionPageResp{
			ID:        questionInfo.ID,
			CreatedAt: questionInfo.CreatedAt.Unix(),
			Title:     questionInfo.Title,
			//UrlTitle:         htmltext.UrlTitle(questionInfo.Title),
			//Description:      htmltext.FetchExcerpt(questionInfo.ParsedText, "...", 240),
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
		}

		questionIDs = append(questionIDs, questionInfo.ID)
		userIDs = append(userIDs, questionInfo.UserID)
		haveEdited, haveAnswered := false, false
		//if checker.IsNotZeroString(questionInfo.LastEditUserID) {
		//	haveEdited = true
		//	userIDs = append(userIDs, questionInfo.LastEditUserID)
		//}
		//if checker.IsNotZeroString(questionInfo.LastAnswerID) {
		//	haveAnswered = true
		//
		//	answerInfo, exist, err := qs.questionDao.AnswerRepo.GetAnswer(ctx, questionInfo.LastAnswerID)
		//	if err == nil && exist {
		//		if answerInfo.LastEditUserID != "0" {
		//			t.LastAnsweredUserID = answerInfo.LastEditUserID
		//		} else {
		//			t.LastAnsweredUserID = answerInfo.UserID
		//		}
		//		t.LastAnsweredAt = answerInfo.CreatedAt
		//		userIDs = append(userIDs, t.LastAnsweredUserID)
		//	}
		//}

		// if order condition is newest or nobody edited or nobody answered, only show question author
		if orderCond == types.QuestionOrderCondNewest || (!haveEdited && !haveAnswered) {
			t.OperationType = types.QuestionPageRespOperationTypeAsked
			t.OperatedAt = questionInfo.CreatedAt.Unix()
			t.Operator = &types.QuestionPageRespOperator{ID: questionInfo.UserID}
		} else {
			// if no one
			if haveEdited {
				t.OperationType = types.QuestionPageRespOperationTypeModified
				t.OperatedAt = questionInfo.UpdatedAt.Unix()
				t.Operator = &types.QuestionPageRespOperator{ID: questionInfo.LastEditUserID}
			}

			if haveAnswered {
				if t.LastAnsweredAt.Unix() > t.OperatedAt {
					t.OperationType = types.QuestionPageRespOperationTypeAnswered
					t.OperatedAt = t.LastAnsweredAt.Unix()
					t.Operator = &types.QuestionPageRespOperator{ID: t.LastAnsweredUserID}
				}
			}
		}
		formattedQuestions = append(formattedQuestions, t)
	}

	//tagsMap, err := TagServicer.BatchGetObjectTag(ctx, questionIDs)
	//if err != nil {
	//	return formattedQuestions, err
	//}
	//userInfoMap, err := UserCommonServicer.BatchUserBasicInfoByID(ctx, userIDs)
	//if err != nil {
	//	return formattedQuestions, err
	//}

	//for _, item := range formattedQuestions {
	//	tags, ok := tagsMap[item.ID]
	//	if ok {
	//		item.Tags = tags
	//	} else {
	//		item.Tags = make([]*types.TagResp, 0)
	//	}
	//	userInfo, ok := userInfoMap[item.Operator.ID]
	//	if ok {
	//		if userInfo != nil {
	//			item.Operator.DisplayName = userInfo.DisplayName
	//			item.Operator.Username = userInfo.Username
	//			item.Operator.Rank = userInfo.Rank
	//			item.Operator.Status = userInfo.Status
	//		}
	//	}
	//
	//}
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

	if questionInfo.Status == types.QuestionStatusDeleted {
		return nil
	}

	questionInfo.Status = types.QuestionStatusDeleted
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
	questionInfo.Status = types.QuestionStatusClosed
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
