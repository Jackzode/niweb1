package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"xorm.io/xorm"
)

type UserDao struct {
	DB    *xorm.Engine
	Cache *redis.Client
}

// NewUserDao new repository
func NewUserDao() *UserDao {
	return &UserDao{
		DB:    handler.Engine,
		Cache: handler.RedisClient,
	}
}

func (ur *UserDao) UpdateUserAvatar(ctx context.Context, uid, url string) (err error) {
	_, err = ur.DB.Context(ctx).ID(uid).Cols("avatar").Update(&types.User{Avatar: url})
	return err
}

// AddUser add user
func (ur *UserDao) AddUser(ctx context.Context, user *types.User) (err error) {

	_, err = ur.DB.Transaction(
		func(session *xorm.Session) (interface{}, error) {
			session = session.Context(ctx)
			userInfo := &types.User{}
			exist, err := session.Where("username = ?", user.Username).Get(userInfo)
			if err != nil {
				return nil, err
			}
			if exist {
				return nil, errors.New(constants.EmailDuplicate)
			}
			_, err = session.Insert(user)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
	return
}

// IncreaseAnswerCount increase answer count
func (ur *UserDao) IncreaseAnswerCount(ctx context.Context, userID string, amount int) (err error) {
	user := &types.User{}
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Incr("answer_count", amount).Update(user)
	if err != nil {
		return err
	}
	return nil
}

// IncreaseQuestionCount increase question count
func (ur *UserDao) IncreaseQuestionCount(ctx context.Context, userID string, amount int) (err error) {
	user := &types.User{}
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Incr("question_count", amount).Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserDao) UpdateQuestionCount(ctx context.Context, userID string, count int64) (err error) {
	user := &types.User{}
	user.QuestionCount = int(count)
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Cols("question_count").Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserDao) UpdateAnswerCount(ctx context.Context, userID string, count int) (err error) {
	user := &types.User{}
	user.AnswerCount = count
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Cols("answer_count").Update(user)
	return err
}

// UpdateLastLoginDate update last login date
func (ur *UserDao) UpdateLastLoginDate(ctx context.Context, userID string) (err error) {
	user := &types.User{LastLoginDate: time.Now()}
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Cols("last_login_date").Update(user)
	return err
}

// UpdateEmailStatus update email status
func (ur *UserDao) UpdateEmailStatus(ctx context.Context, userID string, emailStatus int) error {
	cond := &types.User{MailStatus: emailStatus}
	_, err := ur.DB.Context(ctx).Where("id = ?", userID).Cols("mail_status").Update(cond)
	return err
}

// UpdateNoticeStatus update notice status
func (ur *UserDao) UpdateNoticeStatus(ctx context.Context, userID string, noticeStatus int) error {
	cond := &types.User{NoticeStatus: noticeStatus}
	_, err := ur.DB.Context(ctx).Where("id = ?", userID).Cols("notice_status").Update(cond)
	return err
}

func (ur *UserDao) UpdatePass(ctx context.Context, userID, pass string) error {
	_, err := ur.DB.Context(ctx).Where("id = ?", userID).Cols("pass").Update(&types.User{Pass: pass})
	return err
}

func (ur *UserDao) UpdateEmail(ctx context.Context, userID, email string) (err error) {
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Update(&types.User{EMail: email})
	return
}

func (ur *UserDao) UpdateEmailAndEmailStatus(ctx context.Context, userID, email string, mailStatus int) (err error) {
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Update(&types.User{EMail: email, MailStatus: mailStatus})
	return
}

func (ur *UserDao) UpdateLanguage(ctx context.Context, userID, language string) (err error) {
	_, err = ur.DB.Context(ctx).Where("id = ?", userID).Update(&types.User{Language: language})
	return
}

// UpdateInfo update user info  todo except avatar
func (ur *UserDao) UpdateInfo(ctx context.Context, userInfo *types.User) (err error) {
	_, err = ur.DB.Context(ctx).Where("id = ?", userInfo.ID).
		Cols("username", "display_name", "description", "school", "website", "city_id", "company", "firstname", "lastname", "position", "birthday", "github").Update(userInfo)
	return
}

// GetByUserID get user info by user id
func (ur *UserDao) GetByUserID(ctx context.Context, userID string) (userInfo *types.User, exist bool, err error) {
	userInfo = &types.User{}
	exist, err = ur.DB.Context(ctx).Where("id = ?", userID).Get(userInfo)
	if err != nil {
		return
	}
	//todo  due to plugin
	//err = tryToDecorateUserInfoFromUserCenter(ctx, ur.DB, userInfo)
	//if err != nil {
	//	return nil, false, err
	//}
	return
}

func (ur *UserDao) BatchGetUserByID(ctx context.Context, ids []string) ([]*types.User, error) {
	list := make([]*types.User, 0)
	err := ur.DB.Context(ctx).In("id", ids).Find(&list)
	if err != nil {
		return nil, err
	}
	//tryToDecorateUserListFromUserCenter(ctx, ur.DB, list)
	return list, nil
}

func (ur *UserDao) CheckUsernameValid(ctx context.Context, username string) bool {
	res := ur.Cache.SetNX(ctx, username, 1, 0)
	if res.Err() != nil {
		glog.Slog.Error(res.Err().Error())
	}
	return res.Val()
}

func (ur *UserDao) CheckEmailValid(ctx context.Context, email string) bool {
	res := ur.Cache.SetNX(ctx, email, 1, 0)
	//fmt.Println("setnx---", res.Err(), res.Val())
	if res.Err() != nil {
		glog.Slog.Error(res.Err().Error())
	}
	return res.Val()
}

func (ur *UserDao) CheckEmailExist(ctx context.Context, email string) bool {
	res := ur.Cache.Get(ctx, email)
	if res.Val() == "" || len(res.Val()) == 0 {
		return false
	}
	return true
}

func (ur *UserDao) GetUserInfoByUsername(ctx context.Context, username string) (userInfo *types.User, exist bool, err error) {
	userInfo = &types.User{}
	exist, err = ur.DB.Context(ctx).Where("username = ?", username).Get(userInfo)
	if err != nil {
		return
	}
	//todo
	//err = tryToDecorateUserInfoFromUserCenter(ctx, ur.DB, userInfo)
	//if err != nil {
	//	return nil, false, err
	//}
	return
}

func (ur *UserDao) BatchGetByUsernames(ctx context.Context, usernames []string) ([]*types.User, error) {
	list := make([]*types.User, 0)
	err := ur.DB.Context(ctx).Where("status =?", constants.UserStatusAvailable).In("username", usernames).Find(&list)
	if err != nil {

		return list, err
	}
	//tryToDecorateUserListFromUserCenter(ctx, ur.DB, list)
	return list, nil
}

func (ur *UserDao) GetUserInfoByEmailFromDB(ctx context.Context, email string) (userInfo *types.User, exist bool, err error) {
	userInfo = &types.User{}
	exist, err = ur.DB.Context(ctx).Where("e_mail = ?", email).
		Where("status != ?", constants.UserStatusDeleted).Get(userInfo)
	return
}

func (ur *UserDao) GetUserCount(ctx context.Context) (count int64, err error) {
	session := ur.DB.Context(ctx)
	session.Where("status = ? OR status = ?", constants.UserStatusAvailable, constants.UserStatusSuspended)
	count, err = session.Count(&types.User{})
	if err != nil {
		return count, err
	}
	return count, nil
}

func (ur *UserDao) SearchUserListByName(ctx context.Context, name string, limit int) (userList []*types.User, err error) {
	userList = make([]*types.User, 0)
	session := ur.DB.Context(ctx)
	session.Where("status = ?", constants.UserStatusAvailable)
	session.Where("username LIKE ? OR display_name LIKE ?", strings.ToLower(name)+"%", name+"%")
	session.OrderBy("username ASC, id DESC")
	session.Limit(limit)
	err = session.Find(&userList)
	if err != nil {
		return nil, err
	}
	//todo
	//tryToDecorateUserListFromUserCenter(ctx, ur.DB, userList)
	return
}

func (ur *UserDao) GetCodeContent(ctx context.Context, code string) (info *types.EmailCodeContent, err error) {
	get := ur.Cache.Get(ctx, code)
	if get.Err() != nil {
		glog.Slog.Error(get.Err().Error())
		return nil, get.Err()
	}
	info = &types.EmailCodeContent{}
	err = json.Unmarshal([]byte(get.Val()), info)
	if err != nil {
		glog.Slog.Error(err.Error())
		return nil, err
	}
	return info, nil
}

/*func tryToDecorateUserInfoFromUserCenter(ctx context.Context, db *xorm.Engine, original *types.User) (err error) {
	if original == nil {
		return nil
	}
	uc, ok := plugin.GetUserCenter()
	if !ok {
		return nil
	}

	userInfo := &types.UserExternalLogin{}
	session := db.Context(ctx).Where("user_id = ?", original.ID)
	session.Where("provider = ?", uc.Info().SlugName)
	exist, err := session.Get(userInfo)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	userCenterBasicUserInfo, err := uc.UserInfo(userInfo.ExternalID)
	if err != nil {
		log.Error(err)
		return err
	}

	// In general, usernames should be guaranteed unique by the User Center plugin, so there are no inconsistencies.
	if original.Username != userCenterBasicUserInfo.Username {
		log.Warnf("user %s username is inconsistent with user center", original.ID)
	}
	decorateByUserCenterUser(original, userCenterBasicUserInfo)
	return nil
}

func tryToDecorateUserListFromUserCenter(ctx context.Context, db *xorm.Engine, original []*types.User) {
	uc, ok := plugin.GetUserCenter()
	if !ok {
		return
	}

	ids := make([]string, 0)
	originalUserIDMapping := make(map[string]*types.User, 0)
	for _, user := range original {
		originalUserIDMapping[user.ID] = user
		ids = append(ids, user.ID)
	}

	userExternalLoginList := make([]*types.UserExternalLogin, 0)
	session := db.Context(ctx).Where("provider = ?", uc.Info().SlugName)
	session.In("user_id", ids)
	err := session.Find(&userExternalLoginList)
	if err != nil {
		log.Error(err)
		return
	}

	userExternalIDs := make([]string, 0)
	originalExternalIDMapping := make(map[string]*types.User, 0)
	for _, u := range userExternalLoginList {
		originalExternalIDMapping[u.ExternalID] = originalUserIDMapping[u.UserID]
		userExternalIDs = append(userExternalIDs, u.ExternalID)
	}
	if len(userExternalIDs) == 0 {
		return
	}

	ucUsers, err := uc.UserList(userExternalIDs)
	if err != nil {
		log.Errorf("get user list from user center failed: %v, %v", err, userExternalIDs)
		return
	}

	for _, ucUser := range ucUsers {
		decorateByUserCenterUser(originalExternalIDMapping[ucUser.ExternalID], ucUser)
	}
}

func decorateByUserCenterUser(original *types.User, ucUser *plugin.UserCenterBasicUserInfo) {
	if original == nil || ucUser == nil {
		return
	}
	// In general, usernames should be guaranteed unique by the User Center plugin, so there are no inconsistencies.
	if original.Username != ucUser.Username {
		log.Warnf("user %s username is inconsistent with user center", original.ID)
	}
	if len(ucUser.DisplayName) > 0 {
		original.DisplayName = ucUser.DisplayName
	}
	if len(ucUser.Email) > 0 {
		original.EMail = ucUser.Email
	}
	if len(ucUser.Avatar) > 0 {
		original.Avatar = schema.CustomAvatar(ucUser.Avatar).ToJsonString()
	}
	if len(ucUser.Mobile) > 0 {
		original.Mobile = ucUser.Mobile
	}
	if len(ucUser.Bio) > 0 {
		original.BioHTML = converter.Markdown2HTML(ucUser.Bio) + original.BioHTML
	}

	// If plugin enable rank agent, use rank from user center.
	if plugin.RankAgentEnabled() {
		original.Rank = ucUser.Rank
	}
	if ucUser.Status != plugin.UserStatusAvailable {
		original.Status = int(ucUser.Status)
	}
}


*/
