package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/service/captcha"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

func (us *UserService) EmailLogin(ctx context.Context, req *types.UserEmailLoginReq) (resp *types.UserLoginResp, err error) {

	userInfo, exist, err := us.userDao.GetUserInfoByEmailFromDB(ctx, req.Email)
	if err != nil {
		glog.Slog.Error(err.Error())
		return nil, err
	}
	if !exist || userInfo.Status == constants.UserStatusDeleted {
		glog.Slog.Error(exist, " || userInfo.Status == constants.UserStatusDeleted")
		return nil, errors.New("no exist || user was deleted")
	}
	if !us.verifyPassword(ctx, req.Password, userInfo.Pass) {
		glog.Slog.Error("verify password fail")
		return nil, err
	}
	//更新最近登陆时间
	err = us.userDao.UpdateLastLoginDate(ctx, userInfo.ID)
	if err != nil {
		glog.Slog.Errorf("update last glog.Slogin data failed, err: %v", err)
	}
	resp = &types.UserLoginResp{}
	_ = copier.Copy(resp, userInfo)
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	resp.HavePassword = len(userInfo.Pass) > 0
	resp.RoleID = constants.CommonRole
	resp.Avatar = userInfo.Avatar
	resp.Token, err = utils.CreateToken(userInfo.Username, userInfo.ID, constants.CommonRole)
	if err != nil {
		glog.Slog.Error(err.Error())
		return nil, err
	}
	return resp, nil
}

func (us *UserService) RetrievePassWord(ctx context.Context, req *types.UserRetrievePassWordRequest) error {

	userInfo, has, err := us.userDao.GetUserInfoByEmailFromDB(ctx, req.Email)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	// send email
	data := &types.EmailCodeContent{
		Email:  req.Email,
		UserID: userInfo.ID,
	}
	code := uuid.NewString()
	go captcha.SetCode(ctx, code, utils.JsonObj2String(data), constants.CaptchaExpiration)
	verifyEmailURL := fmt.Sprintf("localhost:8081/users/password-reset?code=%s", code)
	go handler.EmailHandler.Send(req.Email, constants.TitleRetrievePassWord, verifyEmailURL)
	return nil
}

// UpdatePasswordWhenForgot update user password when user forgot password
func (us *UserService) UpdatePasswordWhenForgot(ctx context.Context, req *types.UserRePassWordRequest) (err error) {
	data := &types.EmailCodeContent{}
	//这个content是通过code从缓存里拿到的，里面包含的是用户信息，是用户发修改密码邮件前存的
	err = json.Unmarshal([]byte(req.Content), data)
	if err != nil {
		return err
	}
	//从db中查询用户信息
	userInfo, exist, err := us.userDao.GetUserInfoByEmailFromDB(ctx, data.Email)
	if err != nil || !exist {
		return err
	}
	//加密新密码
	newPass, err := utils.EncryptPassword(req.Pass)
	if err != nil {
		return err
	}
	//更新
	err = us.userDao.UpdatePass(ctx, userInfo.ID, newPass)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UserPassWordVerification(ctx context.Context, uid, oldPass string) bool {
	userInfo, has, err := us.userDao.GetByUserID(ctx, uid)
	if err != nil {
		glog.Slog.Error(err.Error())
		return false
	}
	if !has {
		glog.Slog.Errorf("have not user: %s", uid)
		return false
	}
	isPass := us.verifyPassword(ctx, oldPass, userInfo.Pass)
	if !isPass {
		return false
	}

	return true
}

// UserModifyPassword user modify password
func (us *UserService) UserModifyPassword(ctx context.Context, req *types.UserModifyPasswordReq) error {
	enpass, err := utils.EncryptPassword(req.Pass)
	if err != nil {
		return err
	}
	userInfo, exist, err := us.userDao.GetByUserID(ctx, req.UserID)
	if err != nil || !exist {
		return err
	}

	//再次验证老密码是否正确
	isPass := us.verifyPassword(ctx, req.OldPass, userInfo.Pass)
	if !isPass {
		return errors.New(constants.OldPasswordVerificationFailed)
	}
	//更新数据库密码
	err = us.userDao.UpdatePass(ctx, userInfo.ID, enpass)
	if err != nil {
		return err
	}
	//todo
	return nil
}

func (us *UserService) UpdateInfo(ctx context.Context, req *types.UpdateInfoRequest) (err error) {
	if len(req.Username) > 0 {
		if utils.IsInvalidUsername(req.Username) || utils.IsReservedUsername(req.Username) || utils.IsUsersIgnorePath(req.Username) {
			return errors.New(constants.UsernameInvalid)
		}
		userInfo, exist, err := us.userDao.GetUserInfoByUsername(ctx, req.Username)
		if err != nil {
			glog.Slog.Error(err.Error())
			return err
		}
		if exist && userInfo.ID != req.UserID {
			return errors.New(constants.UsernameDuplicate)
		}
	}

	oldUserInfo, exist, err := us.userDao.GetByUserID(ctx, req.UserID)
	if err != nil {
		glog.Slog.Error(err.Error())
		return err
	}
	if !exist {
		return errors.New(constants.UserNotFound)
	}

	cond := us.formatUserInfoForUpdateInfo(oldUserInfo, req)
	err = us.userDao.UpdateInfo(ctx, cond)
	return err
}

func (us *UserService) formatUserInfoForUpdateInfo(
	oldUserInfo *types.User, req *types.UpdateInfoRequest) *types.User {
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		glog.Slog.Error(err.Error())
	}
	oldUserInfo.Birthday = birthday
	oldUserInfo.Firstname = req.FirstName
	oldUserInfo.Lastname = req.LastName
	oldUserInfo.Github = req.Github
	oldUserInfo.Position = req.Position
	oldUserInfo.DisplayName = req.DisplayName
	oldUserInfo.Username = req.Username
	//todo
	if len(req.Avatar) > 5 {
		oldUserInfo.Avatar = req.Avatar
	}
	oldUserInfo.Company = req.Company
	oldUserInfo.Website = req.Website
	oldUserInfo.CityId = req.CityId
	oldUserInfo.ID = req.UserID
	oldUserInfo.Description = req.Description
	oldUserInfo.School = req.School
	oldUserInfo.Website = req.Website
	oldUserInfo.CityId = req.CityId
	return oldUserInfo
}

// UserUpdateInterface update user interface
func (us *UserService) UserUpdateInterface(ctx context.Context, lang, uid string) (err error) {

	err = us.userDao.UpdateLanguage(ctx, uid, lang)
	if err != nil {
		return
	}
	return nil
}

// UserRegisterByEmail user register
func (us *UserService) UserRegisterByEmail(ctx context.Context, req *types.UserRegisterReq) (resp *types.UserLoginResp, err error) {
	//先查一下数据库是否有这个邮箱地址，有则是重复注册
	//_, has, err := us.userDao.GetUserInfoByEmailFromDB(ctx, req.Email)
	valid := us.userDao.CheckEmailValid(ctx, req.Email)
	//邮箱重复了
	if !valid {
		return nil, errors.New(constants.EmailDuplicate)
	}
	userInfo := &types.User{}
	userInfo.EMail = req.Email
	userInfo.DisplayName = req.Username
	userInfo.Pass, err = utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	userInfo.Username, err = us.MakeUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	userInfo.IPInfo = req.IP
	userInfo.MailStatus = constants.EmailStatusToBeVerified
	userInfo.Status = constants.UserStatusAvailable
	userInfo.LastLoginDate = time.Now()
	userInfo.IPInfo = req.IP
	userInfo.Avatar = constants.DefaultAvatar
	//userInfo.ID是插入mysql生成的
	err = us.userDao.AddUser(ctx, userInfo)
	if err != nil {
		glog.Slog.Error(err.Error())
		return nil, err
	}

	// send email
	data := &types.EmailCodeContent{
		Email:  req.Email,
		UserID: userInfo.ID,
	}
	code := uuid.NewString()
	go captcha.SetCode(ctx, code, utils.JsonObj2String(data), constants.CaptchaExpiration)
	body := fmt.Sprintf("http://localhost:8081/email/verification?code=%s", code)
	go handler.EmailHandler.Send(userInfo.EMail, constants.TitleRegisterByEmail, body)

	// return user info and token
	resp = &types.UserLoginResp{}
	_ = copier.Copy(resp, userInfo)
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	resp.HavePassword = len(userInfo.Pass) > 0
	resp.RoleID = constants.CommonRole
	resp.Avatar = constants.DefaultAvatar
	resp.Token, err = utils.CreateToken(userInfo.Username, userInfo.ID, resp.RoleID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (us *UserService) UserVerifyEmailSend(ctx context.Context, userID string) error {
	userInfo, has, err := us.userDao.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(constants.UserNotFound)
	}

	data := &types.EmailCodeContent{
		Email:  userInfo.EMail,
		UserID: userInfo.ID,
	}
	code := uuid.NewString()
	go captcha.SetCode(ctx, code, utils.JsonObj2String(data), constants.CaptchaExpiration)
	body := fmt.Sprintf("http://localhost:8081/lawyer/user/email/verification?code=%s", code)
	go handler.EmailHandler.Send(userInfo.EMail, constants.TitleVerifyEmail, body)
	return nil
}

func (us *UserService) UserVerifyEmail(ctx context.Context, req *types.UserVerifyEmailReq) (resp *types.UserLoginResp, err error) {
	data := &types.EmailCodeContent{}
	err = utils.FromJsonString2Obj(req.Content, data)
	if err != nil {
		return nil, errors.New(constants.EmailVerifyURLExpired)
	}
	//根据content里的email和uid，查db获取用户的全部信息
	userInfo, has, err := us.userDao.GetUserInfoByEmailFromDB(ctx, data.Email)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(constants.UserNotFound)
	}
	if userInfo.MailStatus == constants.EmailStatusToBeVerified {
		userInfo.MailStatus = constants.EmailStatusAvailable
		//更新用户的邮箱状态为激活状态
		err = us.userDao.UpdateEmailStatus(ctx, userInfo.ID, userInfo.MailStatus)
		if err != nil {
			return nil, err
		}
	}

	AccessToken, err := utils.CreateToken(userInfo.Username, userInfo.ID, 1)
	if err != nil {
		return nil, err
	}

	resp = &types.UserLoginResp{}
	_ = copier.Copy(resp, userInfo)
	resp.Token = AccessToken
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	resp.HavePassword = len(userInfo.Pass) > 0
	resp.Avatar = userInfo.Avatar
	return resp, nil
}

// verifyPassword
// Compare whether the password is correct
func (us *UserService) verifyPassword(ctx context.Context, loginPass, userPass string) bool {
	if len(loginPass) == 0 && len(userPass) == 0 {
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(loginPass))
	return err == nil
}

func (us *UserService) UserChangeEmailSendCode(ctx context.Context, req *types.UserChangeEmailSendCodeReq) (err error) {
	//根据uid查询用户信息
	userInfo, exist, err := us.userDao.GetByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(constants.UserNotFound)
	}
	//校对邮箱状态和密码
	// If user's email already verified, then must verify password first.
	if userInfo.MailStatus == constants.EmailStatusAvailable && !us.verifyPassword(ctx, req.Pass, userInfo.Pass) {

		return errors.New(constants.OldPasswordVerificationFailed)
	}
	//确认下是否是重复的邮箱
	_, exist, err = us.userDao.GetUserInfoByEmailFromDB(ctx, req.Email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(constants.EmailDuplicate)
	}

	data := &types.EmailCodeContent{
		Email:  req.Email,
		UserID: req.UserID,
	}
	code := uuid.NewString()
	go captcha.SetCode(ctx, code, utils.JsonObj2String(data), constants.CaptchaExpiration)

	body := fmt.Sprintf("http://localhost:80/users/confirm-new-email?code=%s", code)
	if userInfo.MailStatus == constants.EmailStatusToBeVerified {

	}
	//给新邮箱发送验证码
	go handler.EmailHandler.Send(req.Email, constants.TitleChangeEmail, body)
	return nil
}

func (us *UserService) UserChangeEmailVerify(ctx context.Context, content string) (resp *types.UserLoginResp, err error) {
	data := &types.EmailCodeContent{}
	err = utils.FromJsonString2Obj(content, data)
	if err != nil {
		return nil, errors.New(constants.EmailVerifyURLExpired)
	}

	_, exist, err := us.userDao.GetUserInfoByEmailFromDB(ctx, data.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New(constants.EmailDuplicate)
	}

	userInfo, exist, err := us.userDao.GetByUserID(ctx, data.UserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(constants.UserNotFound)
	}
	//更新db中的邮箱
	err = us.userDao.UpdateEmailAndEmailStatus(ctx, data.UserID, data.Email, constants.EmailStatusAvailable)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}

	roleID := 1
	resp = &types.UserLoginResp{}
	_ = copier.Copy(resp, userInfo)
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	resp.HavePassword = len(userInfo.Pass) > 0
	resp.Avatar = userInfo.Avatar
	//todo 如何作废之前的token？是个问题
	resp.Token, err = utils.CreateToken(resp.Username, resp.ID, roleID)
	if err != nil {
		return nil, err
	}
	resp.RoleID = roleID
	return resp, nil
}

func (us *UserService) Upload2OSS(filekey string, file io.Reader) (err error) {

	client, err := oss.New("oss-us-west-1.aliyuncs.com", "", "")
	if err != nil {
		return err
	}

	// 获取存储空间。
	bucket, err := client.Bucket("lawyer-niweb1a")
	if err != nil {
		return err
	}
	// 上传文件。
	err = bucket.PutObject(filekey, file)
	return err
}
