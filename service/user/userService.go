package user

import (
	"context"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/dao/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserService struct {
	userDao *user.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		userDao: user.NewUserDao(),
	}
}

func (us *UserService) SearchUserListByName(ctx context.Context, req *types.GetOtherUserInfoByUsernameReq) (
	resp []*types.UserBasicInfo, err error) {
	resp = make([]*types.UserBasicInfo, 0)
	if len(req.Username) == 0 {
		return resp, nil
	}
	//根据username或者display name查db
	userList, err := us.userDao.SearchUserListByName(ctx, req.Username, 5)
	if err != nil {
		return resp, err
	}
	//对检索出来的user拼接一个头像
	for _, u := range userList {
		if req.UserID == u.ID {
			//搜到了自己，就跳过
			continue
		}
		basicInfo := us.FormatUserBasicInfo(ctx, u)
		resp = append(resp, basicInfo)
	}
	return resp, nil
}

func (us *UserService) GetOtherUserInfoByUsername(ctx context.Context, username string) (
	resp *types.GetOtherUserInfoByUsername, err error) {
	//根据username从数据库中获取用户信息
	userInfo, exist, err := us.userDao.GetUserInfoByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, err
	}
	resp = &types.GetOtherUserInfoByUsername{}
	_ = copier.Copy(resp, userInfo)
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	if userInfo.MailStatus == constants.EmailStatusToBeVerified {
		statusMsgShow, ok := constants.UserStatusShowMsg[11]
		if ok {
			resp.StatusMsg = statusMsgShow
		}
	} else {
		statusMsgShow, ok := constants.UserStatusShowMsg[userInfo.Status]
		if ok {
			resp.StatusMsg = statusMsgShow
		}
	}
	return resp, nil
}

func (us *UserService) GetUserBasicInfoByID(ctx context.Context, ID string) (
	userBasicInfo *types.UserBasicInfo, exist bool, err error) {
	userInfo, exist, err := us.userDao.GetByUserID(ctx, ID)
	if err != nil {
		return nil, exist, err
	}
	info := us.FormatUserBasicInfo(ctx, userInfo)
	return info, exist, nil
}

func (us *UserService) GetUserBasicInfoByUserName(ctx context.Context, username string) (*types.UserBasicInfo, bool, error) {
	userInfo, exist, err := us.userDao.GetUserInfoByUsername(ctx, username)
	if err != nil {
		return nil, exist, err
	}
	info := us.FormatUserBasicInfo(ctx, userInfo)
	return info, exist, nil
}

func (us *UserService) BatchGetUserBasicInfoByUserNames(ctx context.Context, usernames []string) (map[string]*types.UserBasicInfo, error) {
	infomap := make(map[string]*types.UserBasicInfo)
	list, err := us.userDao.BatchGetByUsernames(ctx, usernames)
	if err != nil {
		return infomap, err
	}
	for _, user := range list {
		info := us.FormatUserBasicInfo(ctx, user)
		infomap[user.Username] = info
	}
	return infomap, nil
}

func (us *UserService) UpdateAnswerCount(ctx context.Context, userID string, num int) error {
	return us.userDao.UpdateAnswerCount(ctx, userID, num)
}

func (us *UserService) UpdateQuestionCount(ctx context.Context, userID string, num int64) error {
	return us.userDao.UpdateQuestionCount(ctx, userID, num)
}

func (us *UserService) BatchUserBasicInfoByID(ctx context.Context, userIDs []string) (map[string]*types.UserBasicInfo, error) {
	userMap := make(map[string]*types.UserBasicInfo)
	if len(userIDs) == 0 {
		return userMap, nil
	}
	userList, err := us.userDao.BatchGetUserByID(ctx, userIDs)
	if err != nil {
		return userMap, err
	}
	for _, user := range userList {
		info := us.FormatUserBasicInfo(ctx, user)
		userMap[user.ID] = info
	}
	return userMap, nil
}

// FormatUserBasicInfo format user basic info
func (us *UserService) FormatUserBasicInfo(ctx context.Context, userInfo *types.User) *types.UserBasicInfo {
	userBasicInfo := &types.UserBasicInfo{}
	userBasicInfo.ID = userInfo.ID
	userBasicInfo.Username = userInfo.Username
	userBasicInfo.Rank = userInfo.Rank
	userBasicInfo.DisplayName = userInfo.DisplayName
	userBasicInfo.Website = userInfo.Website
	userBasicInfo.CityId = userInfo.CityId
	userBasicInfo.Avatar = userInfo.Avatar
	userBasicInfo.Description = userInfo.Description
	userBasicInfo.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	if userBasicInfo.Status == constants.UserDeleted {
		userBasicInfo.Avatar = ""
		userBasicInfo.DisplayName = "user_delete"
	}
	return userBasicInfo
}

// MakeUsername
// Generate a unique Username based on the displayName
func (us *UserService) MakeUsername(ctx context.Context, displayName string) (username string, err error) {
	//// Chinese processing
	//if has := utils.IsChinese(displayName); has {
	//	str, err := pinyin.New(displayName).Split("").Mode(pinyin.WithoutTone).Convert()
	//	if err != nil {
	//		return "", err
	//	} else {
	//		displayName = str
	//	}
	//}
	//
	//username = strings.ReplaceAll(displayName, " ", "-")
	//username = strings.ToLower(username)
	//suffix := ""

	if utils.IsInvalidUsername(displayName) {
		return "", err
	}
	valid := us.userDao.CheckUsernameValid(ctx, displayName)
	if !valid {
		return "", errors.New("username invalid")
	}
	//todo
	//if utils.IsReservedUsername(username) {
	//	return "", err
	//}

	//for {
	//	_, has, err := us.userDao.GetUserInfoByUsername(ctx, username+suffix)
	//	if err != nil {
	//		return "", err
	//	}
	//	if !has {
	//		break
	//	}
	//	suffix = utils.UsernameSuffix()
	//}
	return displayName, nil
}

func (us *UserService) GetUserInfoByUserID(ctx context.Context, userID string) (
	userInfo *types.User, err error) {

	userInfo, exist, err := us.userDao.GetByUserID(ctx, userID)
	if err != nil {
		glog.Slog.Error(err.Error())
		return nil, err
	}
	if !exist {
		glog.Slog.Error("userID: ", userID, "not found")
		return nil, errors.New(constants.UserNotFound)
	}
	if userInfo.Status == constants.UserStatusDeleted {
		glog.Slog.Error("userID: ", userID, "not found")
		return nil, errors.New(constants.UserDeleted)
	}
	return userInfo, nil

}

func (us *UserService) UploadUserAvatar(ctx *gin.Context, uid string, url string) error {
	return us.userDao.UpdateUserAvatar(ctx, uid, url)
}

/*
// UserRanking get user ranking
func (us *UserService) UserRanking(ctx context.Context) (resp *types.UserRankingResp, err error) {
	limit := 20
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)
	userIDs, userIDExist := make([]string, 0), make(map[string]bool, 0)

	// get most reputation users
	rankStat, rankStatUserIDs, err := us.getActivityUserRankStat(ctx, startTime, endTime, limit, userIDExist)
	if err != nil {
		return nil, err
	}
	userIDs = append(userIDs, rankStatUserIDs...)

	// get most vote users
	voteStat, voteStatUserIDs, err := us.getActivityUserVoteStat(ctx, startTime, endTime, limit, userIDExist)
	if err != nil {
		return nil, err
	}
	userIDs = append(userIDs, voteStatUserIDs...)

	// get all staff members
	userRoleRels, staffUserIDs, err := us.getStaff(ctx, userIDExist)
	if err != nil {
		return nil, err
	}
	userIDs = append(userIDs, staffUserIDs...)

	// get user information
	userInfoMapping, err := us.getUserInfoMapping(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return us.warpStatRankingResp(userInfoMapping, rankStat, voteStat, userRoleRels), nil
}



// UserUnsubscribeNotification user unsubscribe email notification
func (us *UserService) UserUnsubscribeNotification(
	ctx context.Context, req *types.UserUnsubscribeNotificationReq) (err error) {
	data := &types.EmailCodeContent{}
	err = data.FromJSONString(req.Content)
	if err != nil || len(data.UserID) == 0 {
		return errors.New(constants.EmailVerifyURLExpired)
	}

	for _, source := range data.NotificationSources {
		notificationConfig, exist, err := repo.UserNotificationConfigRepo.GetByUserIDAndSource(
			ctx, data.UserID, source)
		if err != nil {
			return err
		}
		if !exist {
			continue
		}
		channels := types.NewNotificationChannelsFormJson(notificationConfig.Channels)
		// unsubscribe email notification
		for _, channel := range channels {
			if channel.Key == constant.EmailChannel {
				channel.Enable = false
			}
		}
		notificationConfig.Channels = channels.ToJsonString()
		if err = repo.UserNotificationConfigRepo.Save(ctx, notificationConfig); err != nil {
			return err
		}
	}
	return nil
}

func (us *UserService) getActivityUserRankStat(ctx context.Context, startTime, endTime time.Time, limit int,
	userIDExist map[string]bool) (rankStat []*types.ActivityUserRankStat, userIDs []string, err error) {
	//if plugin.RankAgentEnabled() {
	//	return make([]*types.ActivityUserRankStat, 0), make([]string, 0), nil
	//}
	rankStat, err = repoCommon.NewActivityRepo().GetUsersWhoHasGainedTheMostReputation(ctx, startTime, endTime, limit)
	if err != nil {
		return nil, nil, err
	}
	for _, stat := range rankStat {
		if stat.Rank <= 0 {
			continue
		}
		//去重
		if userIDExist[stat.UserID] {
			continue
		}
		userIDs = append(userIDs, stat.UserID)
		userIDExist[stat.UserID] = true
	}
	return rankStat, userIDs, nil
}

func (us *UserService) getActivityUserVoteStat(ctx context.Context, startTime, endTime time.Time, limit int,
	userIDExist map[string]bool) (voteStat []*types.ActivityUserVoteStat, userIDs []string, err error) {
	if plugin.RankAgentEnabled() {
		return make([]*types.ActivityUserVoteStat, 0), make([]string, 0), nil
	}
	voteStat, err = repoCommon.NewActivityRepo().GetUsersWhoHasVoteMost(ctx, startTime, endTime, limit)
	if err != nil {
		return nil, nil, err
	}
	for _, stat := range voteStat {
		if stat.VoteCount <= 0 {
			continue
		}
		if userIDExist[stat.UserID] {
			continue
		}
		userIDs = append(userIDs, stat.UserID)
		userIDExist[stat.UserID] = true
	}
	return voteStat, userIDs, nil
}

func (us *UserService) getStaff(ctx context.Context, userIDExist map[string]bool) (
	userRoleRels []*types.UserRoleRel, userIDs []string, err error) {
	userRoleRels, err = UserRoleRelServicer.GetUserByRoleID(ctx, []int{RoleAdminID, RoleModeratorID})
	if err != nil {
		return nil, nil, err
	}
	for _, rel := range userRoleRels {
		if userIDExist[rel.UserID] {
			continue
		}
		userIDs = append(userIDs, rel.UserID)
		userIDExist[rel.UserID] = true
	}
	return userRoleRels, userIDs, nil
}

func (us *UserService) getUserInfoMapping(ctx context.Context, userIDs []string) (
	userInfoMapping map[string]*types.User, err error) {
	userInfoMapping = make(map[string]*types.User, 0)
	if len(userIDs) == 0 {
		return userInfoMapping, nil
	}
	userInfoList, err := us.userDao.BatchGetByID(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	avatarMapping := types.FormatListAvatar(userInfoList)
	for _, user := range userInfoList {
		user.Avatar = avatarMapping[user.ID].GetURL()
		userInfoMapping[user.ID] = user
	}
	return userInfoMapping, nil
}



func (us *UserService) warpStatRankingResp(
	userInfoMapping map[string]*types.User,
	rankStat []*types.ActivityUserRankStat,
	voteStat []*types.ActivityUserVoteStat,
	userRoleRels []*types.UserRoleRel) (resp *types.UserRankingResp) {
	resp = &types.UserRankingResp{
		UsersWithTheMostReputation: make([]*types.UserRankingSimpleInfo, 0),
		UsersWithTheMostVote:       make([]*types.UserRankingSimpleInfo, 0),
		Staffs:                     make([]*types.UserRankingSimpleInfo, 0),
	}
	for _, stat := range rankStat {
		if stat.Rank <= 0 {
			continue
		}
		if userInfo := userInfoMapping[stat.UserID]; userInfo != nil && userInfo.Status != types.UserStatusDeleted {
			resp.UsersWithTheMostReputation = append(resp.UsersWithTheMostReputation, &types.UserRankingSimpleInfo{
				Username:    userInfo.Username,
				Rank:        stat.Rank,
				DisplayName: userInfo.DisplayName,
				Avatar:      userInfo.Avatar,
			})
		}
	}
	for _, stat := range voteStat {
		if stat.VoteCount <= 0 {
			continue
		}
		if userInfo := userInfoMapping[stat.UserID]; userInfo != nil && userInfo.Status != types.UserStatusDeleted {
			resp.UsersWithTheMostVote = append(resp.UsersWithTheMostVote, &types.UserRankingSimpleInfo{
				Username:    userInfo.Username,
				VoteCount:   stat.VoteCount,
				DisplayName: userInfo.DisplayName,
				Avatar:      userInfo.Avatar,
			})
		}
	}
	for _, rel := range userRoleRels {
		if userInfo := userInfoMapping[rel.UserID]; userInfo != nil && userInfo.Status != types.UserStatusDeleted {
			resp.Staffs = append(resp.Staffs, &types.UserRankingSimpleInfo{
				Username:    userInfo.Username,
				Rank:        userInfo.Rank,
				DisplayName: userInfo.DisplayName,
				Avatar:      userInfo.Avatar,
			})
		}
	}
	return resp
}

func (us *UserService) CacheLoginUserInfo(ctx context.Context, userID string, userStatus, emailStatus int, externalID string) (
	accessToken string, userCacheInfo *types.UserCacheInfo, err error) {

	roleID, err := UserRoleRelServicer.GetUserRole(ctx, userID)
	if err != nil {
		glog.Slog.Error(err)
	}

	userCacheInfo = &types.UserCacheInfo{
		UserID:      userID,
		EmailStatus: emailStatus,
		UserStatus:  userStatus,
		RoleID:      roleID,
		ExternalID:  externalID,
	}

	accessToken, err = AuthServicer.SetUserCacheInfo(ctx, userCacheInfo)
	if err != nil {
		return "", nil, err
	}
	return accessToken, userCacheInfo, nil
}
*/
