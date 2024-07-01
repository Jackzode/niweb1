package comment

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/dao/comment"
	"github.com/jinzhu/copier"
	"time"
)

type CommentService struct {
	cd *comment.CommentDao
}

// NewCommentCommonService new comment service
func NewCommentCommonService() *CommentService {
	return &CommentService{
		cd: comment.NewCommentDao(),
	}
}

// GetCommentWithPage get comment list page
func (cs *CommentService) GetCommentWithPage(ctx context.Context, req *types.GetCommentWithPageReq) (
	resp []*types.GetCommentResp, total int64, err error) {
	dto := &types.CommentQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		ObjectID:  req.ObjectID,
		QueryCond: req.QueryCond,
	}
	commentList, total, err := cs.cd.GetCommentPage(ctx, dto)
	if err != nil {
		return nil, 0, err
	}
	resp = make([]*types.GetCommentResp, 0)
	for _, each := range commentList {
		commentResp, err := cs.convertCommentEntity2Resp(ctx, req, each)
		if err != nil {
			return nil, 0, err
		}
		resp = append(resp, commentResp)
	}

	// if user request the specific each, add it if not exist.
	if len(req.CommentID) > 0 {
		commentExist := false
		for _, t := range resp {
			if t.CommentID == req.CommentID {
				commentExist = true
				break
			}
		}
		if !commentExist {
			eachComment, exist, err := cs.cd.GetCommentById(ctx, req.CommentID)
			if err != nil {
				return nil, 0, err
			}
			if exist && eachComment.ObjectID == req.ObjectID {
				commentResp, err := cs.convertCommentEntity2Resp(ctx, req, eachComment)
				if err != nil {
					return nil, 0, err
				}
				resp = append(resp, commentResp)
			}
		}
	}
	return resp, total, nil
}

// GetComment get comment one
func (cs *CommentService) GetComment(ctx context.Context, req *types.GetCommentReq) (resp *types.GetCommentResp, err error) {
	comments, exist, err := cs.cd.GetCommentById(ctx, req.ID)
	if err != nil {
		return
	}
	if !exist {
		return nil, errors.New(constants.CommentNotFound)
	}

	resp = &types.GetCommentResp{
		CommentID:      comments.ID,
		CreatedAt:      comments.CreatedAt.Unix(),
		UserID:         comments.UserID,
		ReplyUserID:    comments.GetReplyUserID(),
		ReplyCommentID: comments.GetReplyCommentID(),
		ObjectID:       comments.ObjectID,
		VoteCount:      comments.VoteCount,
		OriginalText:   comments.OriginalText,
		ParsedText:     comments.ParsedText,
	}

	// get comments user info
	/*if len(resp.UserID) > 0 {
		commentUser, exist, err := cs.us.GetUserBasicInfoByID(ctx, resp.UserID)
		if err != nil {
			return nil, err
		}
		if exist {
			resp.Username = commentUser.Username
			resp.UserDisplayName = commentUser.DisplayName
			resp.UserAvatar = commentUser.Avatar
			resp.UserStatus = commentUser.Status
		}
	}

	// get reply user info
	if len(resp.ReplyUserID) > 0 {
		replyUser, exist, err := UserCommonServicer.GetUserBasicInfoByID(ctx, resp.ReplyUserID)
		if err != nil {
			return nil, err
		}
		if exist {
			resp.ReplyUsername = replyUser.Username
			resp.ReplyUserDisplayName = replyUser.DisplayName
			resp.ReplyUserStatus = replyUser.Status
		}
	}

	*/

	// check if current user vote this comments
	//resp.IsVote = cs.checkIsVote(ctx, req.UserID, resp.CommentID)

	//resp.MemberActions = permission.GetCommentPermission(ctx, req.UserID, resp.UserID, comments.CreatedAt, req.CanEdit, req.CanDelete)
	return resp, nil
}

// RemoveComment delete comment
func (cs *CommentService) RemoveComment(ctx context.Context, req *types.RemoveCommentReq) (err error) {
	return cs.cd.RemoveComment(ctx, req.CommentID)
}

// AddComment add comment
func (cs *CommentService) AddComment(ctx context.Context, req *types.AddCommentReq) (
	resp *types.GetCommentResp, err error) {
	oneComment := &types.Comment{}
	_ = copier.Copy(oneComment, req)
	oneComment.Status = constants.CommentStatusAvailable

	// add question id
	//objInfo, err := ObjServicer.GetInfo(ctx, req.ObjectID)
	//if err != nil {
	//	return nil, err
	//}
	//objInfo.ObjectID = utils.DeShortID(objInfo.ObjectID)
	//objInfo.QuestionID = utils.DeShortID(objInfo.QuestionID)
	//objInfo.AnswerID = utils.DeShortID(objInfo.AnswerID)
	//if objInfo.ObjectType == constants.QuestionObjectType || objInfo.ObjectType == constants.AnswerObjectType {
	//	oneComment.QuestionID = objInfo.QuestionID
	//}

	if len(req.ReplyCommentID) > 0 {
		replyComment, exist, err := cs.cd.GetCommentById(ctx, req.ReplyCommentID)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errors.New(constants.CommentNotFound)
		}
		oneComment.SetReplyUserID(replyComment.UserID)
		oneComment.SetReplyCommentID(replyComment.ID)
	} else {
		oneComment.SetReplyUserID("")
		oneComment.SetReplyCommentID("")
	}

	err = cs.cd.AddComment(ctx, oneComment)
	if err != nil {
		return nil, err
	}

	resp = &types.GetCommentResp{}
	resp.SetFromComment(oneComment)
	//resp.MemberActions = permission.GetCommentPermission(ctx, req.UserID, resp.UserID, time.Now(), req.CanEdit, req.CanDelete)

	//commentResp, err := cs.addCommentNotification(ctx, req, resp, oneComment, objInfo)
	//if err != nil {
	//	return commentResp, err
	//}

	// todo get user info
	//userInfo, exist, err := UserCommonServicer.GetUserBasicInfoByID(ctx, resp.UserID)
	//if err != nil {
	//	return nil, err
	//}
	//if exist {
	//	resp.Username = userInfo.Username
	//	resp.UserDisplayName = userInfo.DisplayName
	//	resp.UserAvatar = userInfo.Avatar
	//	resp.UserStatus = userInfo.Status
	//}

	//activityMsg := &types.ActivityMsg{
	//	UserID:           oneComment.UserID,
	//	ObjectID:         oneComment.ID,
	//	OriginalObjectID: req.ObjectID,
	//	ActivityTypeKey:  constant.ActQuestionCommented,
	//}
	//switch objInfo.ObjectType {
	//case constant.QuestionObjectType:
	//	activityMsg.ActivityTypeKey = constant.ActQuestionCommented
	//case constant.AnswerObjectType:
	//	activityMsg.ActivityTypeKey = constant.ActAnswerCommented
	//}
	//ActivityQueueServicer.Send(ctx, activityMsg)
	return resp, nil
}

// UpdateComment update comment
func (cs *CommentService) UpdateComment(ctx context.Context, req *types.UpdateCommentReq) (
	resp *types.UpdateCommentResp, err error) {
	old, exist, err := cs.cd.GetCommentById(ctx, req.CommentID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(constants.CommentNotFound)
	}
	// user can't edit the comment that was posted by others except admin
	if !req.IsAdmin && req.UserID != old.UserID {
		return nil, errors.New(constants.CommentNotFound)
	}

	// user can edit the comment that was posted by himself before deadline.
	// admin can edit it at any time
	if !req.IsAdmin && (time.Now().After(old.CreatedAt.Add(constants.CommentEditDeadline))) {
		return nil, errors.New(constants.CommentCannotEditAfterDeadline)
	}

	if err = cs.cd.UpdateCommentContent(ctx, old.ID, req.OriginalText, req.ParsedText); err != nil {
		return nil, err
	}
	resp = &types.UpdateCommentResp{
		CommentID:    old.ID,
		OriginalText: req.OriginalText,
		ParsedText:   req.ParsedText,
	}
	return resp, nil
}

func (cs *CommentService) convertCommentEntity2Resp(ctx context.Context, req *types.GetCommentWithPageReq,
	comment *types.Comment) (commentResp *types.GetCommentResp, err error) {
	commentResp = &types.GetCommentResp{
		CommentID:      comment.ID,
		CreatedAt:      comment.CreatedAt.Unix(),
		UserID:         comment.UserID,
		ReplyUserID:    comment.GetReplyUserID(),
		ReplyCommentID: comment.GetReplyCommentID(),
		ObjectID:       comment.ObjectID,
		VoteCount:      comment.VoteCount,
		OriginalText:   comment.OriginalText,
		ParsedText:     comment.ParsedText,
	}

	// get comment user info
	//if len(commentResp.UserID) > 0 {
	//	commentUser, exist, err := UserCommonServicer.GetUserBasicInfoByID(ctx, commentResp.UserID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if exist {
	//		commentResp.Username = commentUser.Username
	//		commentResp.UserDisplayName = commentUser.DisplayName
	//		commentResp.UserAvatar = commentUser.Avatar
	//		commentResp.UserStatus = commentUser.Status
	//	}
	//}

	// get reply user info
	//if len(commentResp.ReplyUserID) > 0 {
	//	replyUser, exist, err := UserCommonServicer.GetUserBasicInfoByID(ctx, commentResp.ReplyUserID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if exist {
	//		commentResp.ReplyUsername = replyUser.Username
	//		commentResp.ReplyUserDisplayName = replyUser.DisplayName
	//		commentResp.ReplyUserStatus = replyUser.Status
	//	}
	//}

	// check if current user vote this comment
	//commentResp.IsVote = cs.checkIsVote(ctx, req.UserID, commentResp.CommentID)

	//commentResp.MemberActions = permission.GetCommentPermission(ctx, req.UserID, commentResp.UserID, comment.CreatedAt, req.CanEdit, req.CanDelete)
	return commentResp, nil
}

//func (cs *CommentService) checkIsVote(ctx context.Context, userID, commentID string) (isVote bool) {
//	status := repo.VoteRepo.GetVoteStatus(ctx, commentID, userID)
//	return len(status) > 0
//}

// GetCommentPersonalWithPage get personal comment list page
func (cs *CommentService) GetCommentPersonalWithPage(ctx context.Context, req *types.GetCommentPersonalWithPageReq) (
	resp []*types.GetCommentPersonalWithPageResp, total int64, err error) {
	//if len(req.Username) > 0 {
	//	userInfo, exist, err := UserCommonServicer.GetUserBasicInfoByUserName(ctx, req.Username)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if !exist {
	//		return nil, errors.New(constants.UserNotFound)
	//	}
	//	req.UserID = userInfo.ID
	//}
	if len(req.UserID) == 0 {
		return nil, 0, errors.New(constants.UserNotFound)
	}

	dto := &types.CommentQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		UserID:    req.UserID,
		QueryCond: "created_at",
	}
	commentList, total, err := cs.cd.GetCommentPage(ctx, dto)
	if err != nil {
		return nil, 0, err
	}
	resp = make([]*types.GetCommentPersonalWithPageResp, 0)
	for _, each := range commentList {
		commentResp := &types.GetCommentPersonalWithPageResp{
			CommentID: each.ID,
			CreatedAt: each.CreatedAt.Unix(),
			ObjectID:  each.ObjectID,
			Content:   each.ParsedText, // todo trim
		}
		//if len(each.ObjectID) > 0 {
		//	objInfo, err := ObjServicer.GetInfo(ctx, each.ObjectID)
		//	if err != nil {
		//		log.Error(err)
		//	} else {
		//		commentResp.ObjectType = objInfo.ObjectType
		//		commentResp.Title = objInfo.Title
		//		commentResp.UrlTitle = htmltext.UrlTitle(objInfo.Title)
		//		commentResp.QuestionID = objInfo.QuestionID
		//		commentResp.AnswerID = objInfo.AnswerID
		//		if objInfo.QuestionStatus == types.QuestionStatusDeleted {
		//			commentResp.Title = "Deleted question"
		//		}
		//	}
		//}
		resp = append(resp, commentResp)
	}
	return resp, total, nil
}

/*
func (cs *CommentService) notificationQuestionComment(ctx context.Context, questionUserID,
	questionID, questionTitle, commentID, commentUserID, commentSummary string) {
	if questionUserID == commentUserID {
		return
	}
	// send internal notification
	msg := &types.NotificationMsg{
		ReceiverUserID: questionUserID,
		TriggerUserID:  commentUserID,
		Type:           types.NotificationTypeInbox,
		ObjectID:       commentID,
	}
	msg.ObjectType = constant.CommentObjectType
	msg.NotificationAction = constant.NotificationCommentQuestion
	NotificationQueueService.Send(ctx, msg)

	// send external notification
	receiverUserInfo, exist, err := repo.UserRepo.GetByUserID(ctx, questionUserID)
	if err != nil {
		log.Error(err)
		return
	}
	if !exist {
		log.Warnf("user %s not found", questionUserID)
		return
	}

	externalNotificationMsg := &types.ExternalNotificationMsg{
		ReceiverUserID: receiverUserInfo.ID,
		ReceiverEmail:  receiverUserInfo.EMail,
		ReceiverLang:   receiverUserInfo.Language,
	}
	rawData := &types.NewCommentTemplateRawData{
		QuestionTitle:   questionTitle,
		QuestionID:      questionID,
		CommentID:       commentID,
		CommentSummary:  commentSummary,
		UnsubscribeCode: token.GenerateToken(),
	}
	commentUser, _, _ := UserCommonServicer.GetUserBasicInfoByID(ctx, commentUserID)
	if commentUser != nil {
		rawData.CommentUserDisplayName = commentUser.DisplayName
	}
	externalNotificationMsg.NewCommentTemplateRawData = rawData
	ExternalNotificationQueueService.Send(ctx, externalNotificationMsg)
}

func (cs *CommentService) notificationAnswerComment(ctx context.Context,
	questionID, questionTitle, answerID, answerUserID, commentID, commentUserID, commentSummary string) {
	if answerUserID == commentUserID {
		return
	}

	// Send internal notification.
	msg := &types.NotificationMsg{
		ReceiverUserID: answerUserID,
		TriggerUserID:  commentUserID,
		Type:           types.NotificationTypeInbox,
		ObjectID:       commentID,
	}
	msg.ObjectType = constant.CommentObjectType
	msg.NotificationAction = constant.NotificationCommentAnswer
	NotificationQueueService.Send(ctx, msg)

	// Send external notification.
	receiverUserInfo, exist, err := repo.UserRepo.GetByUserID(ctx, answerUserID)
	if err != nil {
		log.Error(err)
		return
	}
	if !exist {
		log.Warnf("user %s not found", answerUserID)
		return
	}
	externalNotificationMsg := &types.ExternalNotificationMsg{
		ReceiverUserID: receiverUserInfo.ID,
		ReceiverEmail:  receiverUserInfo.EMail,
		ReceiverLang:   receiverUserInfo.Language,
	}
	rawData := &types.NewCommentTemplateRawData{
		QuestionTitle:   questionTitle,
		QuestionID:      questionID,
		AnswerID:        answerID,
		CommentID:       commentID,
		CommentSummary:  commentSummary,
		UnsubscribeCode: token.GenerateToken(),
	}
	commentUser, _, _ := UserCommonServicer.GetUserBasicInfoByID(ctx, commentUserID)
	if commentUser != nil {
		rawData.CommentUserDisplayName = commentUser.DisplayName
	}
	externalNotificationMsg.NewCommentTemplateRawData = rawData
	ExternalNotificationQueueService.Send(ctx, externalNotificationMsg)
}

func (cs *CommentService) notificationCommentReply(ctx context.Context, replyUserID, commentID, commentUserID string) {
	msg := &types.NotificationMsg{
		ReceiverUserID: replyUserID,
		TriggerUserID:  commentUserID,
		Type:           types.NotificationTypeInbox,
		ObjectID:       commentID,
	}
	msg.ObjectType = constant.CommentObjectType
	msg.NotificationAction = constant.NotificationReplyToYou
	NotificationQueueService.Send(ctx, msg)
}

func (cs *CommentService) notificationMention(
	ctx context.Context, mentionUsernameList []string, commentID, commentUserID string,
	alreadyNotifiedUserID map[string]bool) (alreadyNotifiedUserIDs []string) {
	for _, username := range mentionUsernameList {
		userInfo, exist, err := UserCommonServicer.GetUserBasicInfoByUserName(ctx, username)
		if err != nil {
			log.Error(err)
			continue
		}
		if exist && !alreadyNotifiedUserID[userInfo.ID] {
			msg := &types.NotificationMsg{
				ReceiverUserID: userInfo.ID,
				TriggerUserID:  commentUserID,
				Type:           types.NotificationTypeInbox,
				ObjectID:       commentID,
			}
			msg.ObjectType = constant.CommentObjectType
			msg.NotificationAction = constant.NotificationMentionYou
			NotificationQueueService.Send(ctx, msg)
			alreadyNotifiedUserIDs = append(alreadyNotifiedUserIDs, userInfo.ID)
		}
	}
	return alreadyNotifiedUserIDs
}



func (cs *CommentService) addCommentNotification(
	ctx context.Context, req *types.AddCommentReq, resp *types.GetCommentResp,
	comment *types.Comment, objInfo *types.SimpleObjectInfo) (*types.GetCommentResp, error) {
	// The priority of the notification
	// 1. reply to user
	// 2. comment mention to user
	// 3. answer or question was commented
	alreadyNotifiedUserID := make(map[string]bool)

	// get reply user info
	if len(resp.ReplyUserID) > 0 && resp.ReplyUserID != req.UserID {
		replyUser, exist, err := UserCommonServicer.GetUserBasicInfoByID(ctx, resp.ReplyUserID)
		if err != nil {
			return nil, err
		}
		if exist {
			resp.ReplyUsername = replyUser.Username
			resp.ReplyUserDisplayName = replyUser.DisplayName
			resp.ReplyUserStatus = replyUser.Status
		}
		cs.notificationCommentReply(ctx, replyUser.ID, comment.ID, req.UserID)
		alreadyNotifiedUserID[replyUser.ID] = true
		return nil, nil
	}

	if len(req.MentionUsernameList) > 0 {
		alreadyNotifiedUserIDs := cs.notificationMention(
			ctx, req.MentionUsernameList, comment.ID, req.UserID, alreadyNotifiedUserID)
		for _, userID := range alreadyNotifiedUserIDs {
			alreadyNotifiedUserID[userID] = true
		}
		return nil, nil
	}

	if objInfo.ObjectType == constant.QuestionObjectType && !alreadyNotifiedUserID[objInfo.ObjectCreatorUserID] {
		cs.notificationQuestionComment(ctx, objInfo.ObjectCreatorUserID,
			objInfo.QuestionID, objInfo.Title, comment.ID, req.UserID, comment.OriginalText)
	} else if objInfo.ObjectType == constant.AnswerObjectType && !alreadyNotifiedUserID[objInfo.ObjectCreatorUserID] {
		cs.notificationAnswerComment(ctx, objInfo.QuestionID, objInfo.Title, objInfo.AnswerID,
			objInfo.ObjectCreatorUserID, comment.ID, req.UserID, comment.OriginalText)
	}
	return nil, nil
}

*/
