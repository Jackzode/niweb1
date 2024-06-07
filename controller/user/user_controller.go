package controller

import (
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/commons/types"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/Jackzode/painting/controller"
	"github.com/Jackzode/painting/service/captcha"
	"github.com/Jackzode/painting/service/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// UserController user controller, no need login
type UserController struct {
	userService *user.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: user.NewUserService(),
	}
}

func (uc *UserController) UserRegisterByEmail(ctx *gin.Context) {
	req := &types.UserRegisterReq{}
	if !controller.BindAndCheck(ctx, req) {
		return
	}
	req.IP = ctx.ClientIP()
	//对比验证码是否正确
	//captchaPass, err := captcha.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode)
	//fmt.Println(captchaPass, err)
	//if err != nil || !captchaPass {
	//	controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
	//	return
	//}
	//核心逻辑
	resp, err := uc.userService.UserRegisterByEmail(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, err.Error(), nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) UserEmailLogin(ctx *gin.Context) {
	req := &types.UserEmailLoginReq{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//验证码是否正确
	captchaPass, err := captcha.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode)
	if err != nil || !captchaPass {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}
	resp, err := user.NewUserService().EmailLogin(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
		return
	}

	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) GetUserInfoByUserID(ctx *gin.Context) {

	uid := utils.GetUidFromTokenByCtx(ctx)
	userInfo, err := uc.userService.GetUserInfoByUserID(ctx, uid)
	if err != nil {
		glog.Slog.Error(err)
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
		return
	}
	resp := &types.UserLoginResp{}
	_ = copier.Copy(resp, userInfo)
	resp.CreatedAt = userInfo.CreatedAt.Unix()
	resp.LastLoginDate = userInfo.LastLoginDate.Unix()
	resp.Status = utils.ConvertUserStatus(userInfo.Status, userInfo.MailStatus)
	resp.HavePassword = len(userInfo.Pass) > 0
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) GetOtherUserInfoByUsername(ctx *gin.Context) {
	req := &types.GetOtherUserInfoByUsernameReq{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//核心逻辑
	resp, err := uc.userService.GetOtherUserInfoByUsername(ctx, req.Username)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) RetrievePassWord(ctx *gin.Context) {
	req := &types.UserRetrievePassWordRequest{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//校对验证码
	captchaPass, err := captcha.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode)
	if err != nil || !captchaPass {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}
	//core logic
	err = uc.userService.RetrievePassWord(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserReplacePassWord(ctx *gin.Context) {
	req := &types.UserRePassWordRequest{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//这个code是/password/reset接口生成的，里面存的是email和uid
	req.Content, _ = captcha.GetContentByCaptchaCode(ctx, req.Code)
	if len(req.Content) == 0 {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.EmailVerifyURLExpired, nil)
		return
	}
	//更新db中的密码
	err := uc.userService.UpdatePasswordWhenForgot(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserVerifyEmail(ctx *gin.Context) {
	req := &types.UserVerifyEmailReq{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//VerifyEmailByCode 根据code从缓存中获取content,包含email和uid
	req.Content, _ = captcha.GetContentByCaptchaCode(ctx, req.Code)
	if len(req.Content) == 0 {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.EmailVerifyURLExpired, nil)
		return
	}
	//验证邮箱
	resp, err := uc.userService.UserVerifyEmail(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) UserVerifyEmailSend(ctx *gin.Context) {
	req := &types.UserVerifyEmailReq{}
	//todo  content need verify
	if !controller.BindAndCheck(ctx, req) {
		return
	}
	uid := utils.GetUidFromTokenByCtx(ctx)
	captchaPass, err := captcha.VerifyCaptcha(ctx, req.Code, req.Content)
	if err != nil || !captchaPass {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}

	err = uc.userService.UserVerifyEmailSend(ctx, uid)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserModifyPassWord(ctx *gin.Context) {
	req := &types.UserModifyPasswordReq{}
	fmt.Println("req ", req)
	if controller.BindAndCheck(ctx, req) {
		return
	}
	uid := utils.GetUidFromTokenByCtx(ctx)
	req.UserID = uid
	//校对验证码
	captchaPass, err := captcha.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode)
	if err != nil || !captchaPass {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}

	//验证用户老密码是否正确
	oldPassVerification := uc.userService.UserPassWordVerification(ctx, req.UserID, req.OldPass)
	if !oldPassVerification {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.OldPasswordVerificationFailed, nil)
		return
	}

	//修改密码时新密码和老密码不能一样
	if req.OldPass == req.Pass {
		controller.HandleResponse(ctx, constants.ParamInvalid, constants.NewPasswordSameAsPreviousSetting, nil)
		return
	}
	err = uc.userService.UserModifyPassword(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserUpdateInfo(ctx *gin.Context) {
	req := &types.UpdateInfoRequest{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	//从token里获取用户信息
	req.UserID = utils.GetUidFromTokenByCtx(ctx)
	err := uc.userService.UpdateInfo(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserUpdateInterfaceLang(ctx *gin.Context) {
	//req := &types.UpdateUserInterfaceRequest{}
	//if controller.BindAndCheck(ctx, req) {
	//	return
	//}
	////req.UserId = middleware.GetLoginUserIDFromContext(ctx)
	//req.UserId = utils.GetUidFromTokenByCtx(ctx)
	//if !translator.CheckLanguageIsValid(req.Language) {
	//	controller.HandleResponse(ctx, errors.New(constants.LangNotFound), nil)
	//	return
	//}
	////根据uid更新用户user表的language字段
	//err := service.UserServicer.UserUpdateInterface(ctx, req)
	//controller.HandleResponse(ctx, err, nil)
}

func (uc *UserController) ActionRecord(ctx *gin.Context) {
	//req := &types.ActionRecordReq{}
	//if controller.BindAndCheck(ctx, req) {
	//	return
	//}
	//uid := utils.GetUidFromTokenByCtx(ctx)
	//req.UserID = uid
	//req.IP = ctx.ClientIP()
	//resp := &types.ActionRecordResp{}
	////role id 是2和3，就是管理员， 管理员不需要验证
	////service.CaptchaServicer.ActionRecordAdd(ctx, req.Action, req.IP)
	//unit := service.CaptchaServicer.GetActionRecordUnit(ctx, req)
	////对于当前action是否需要验证码
	//verificationResult := service.CaptchaServicer.ValidationStrategy(ctx, unit, req.Action)
	////需要验证码
	//var err error
	//if verificationResult {
	//	resp.CaptchaID, resp.CaptchaImg, err = service.CaptchaServicer.GenerateCaptcha(ctx)
	//	resp.Verify = true
	//}
	//controller.HandleResponse(ctx, err, resp)

}

func (uc *UserController) UserRegisterCaptcha(ctx *gin.Context) {
	resp := &types.ActionRecordResp{}
	key, base64, err := captcha.GenerateCaptchaAndSave(ctx)
	if err != nil {
		glog.Slog.Error(err.Error())
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}
	resp.Verify = true
	resp.CaptchaID = key
	resp.CaptchaImg = base64
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}

func (uc *UserController) GetUserNotificationConfig(ctx *gin.Context) {
	////userID := middleware.GetLoginUserIDFromContext(ctx)
	//userID := utils.GetUidFromTokenByCtx(ctx)
	//resp, err := service.UserNotificationConfigService.GetUserNotificationConfig(ctx, userID)
	//controller.HandleResponse(ctx, err, resp)
}

// todo @Router /answer/api/v1/user/notification/config [put]
func (uc *UserController) UpdateUserNotificationConfig(ctx *gin.Context) {
	//req := &types.UpdateUserNotificationConfigReq{}
	//if controller.BindAndCheck(ctx, req) {
	//	return
	//}
	//
	////req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//req.UserID = utils.GetUidFromTokenByCtx(ctx)
	//err := service.UserNotificationConfigService.UpdateUserNotificationConfig(ctx, req)
	//controller.HandleResponse(ctx, err, nil)
}

func (uc *UserController) UserChangeEmailSendCode(ctx *gin.Context) {
	req := &types.UserChangeEmailSendCodeReq{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = utils.GetUidFromTokenByCtx(ctx)
	// If the user is not logged in, the api cannot be used.
	// If the user email is not verified, that also can use this api to modify the email.
	if len(req.UserID) == 0 {
		controller.HandleResponse(ctx, constants.StatusErrCode, constants.UnauthorizedError, nil)
		return
	}

	//校对验证码
	captchaPass, err := captcha.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaCode)
	//记录本次修改用户信息的操作
	if err != nil || !captchaPass {
		controller.HandleResponse(ctx, constants.CaptchaFailedCode, constants.CaptchaVerificationFailed, nil)
		return
	}
	//核心逻辑
	err = uc.userService.UserChangeEmailSendCode(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.InternalErrMsg, nil)
		return
	}
	//删除这个操作记录
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, nil)
}

func (uc *UserController) UserChangeEmailVerify(ctx *gin.Context) {
	//req := &types.UserChangeEmailVerifyReq{}
	//if controller.BindAndCheck(ctx, req) {
	//	return
	//}
	//req.Content = service.EmailServicer.VerifyEmailByCode(ctx, req.Code)
	//if len(req.Content) == 0 {
	//	controller.HandleResponse(ctx, errors.New(constants.EmailVerifyURLExpired), nil)
	//	return
	//}
	////核心逻辑
	//resp, err := service.UserServicer.UserChangeEmailVerify(ctx, req.Content)
	//service.CaptchaServicer.ActionRecordDel(ctx, entity.CaptchaActionEmail, ctx.ClientIP())
	//controller.HandleResponse(ctx, err, resp)
}

//func (uc *UserController) UserRanking(ctx *gin.Context) {
//	resp, err := service.UserServicer.UserRanking(ctx)
//	controller.HandleResponse(ctx, err, resp)
//}

//func (uc *UserController) UserUnsubscribeNotification(ctx *gin.Context) {
//	req := &types.UserUnsubscribeNotificationReq{}
//	if controller.BindAndCheck(ctx, req) {
//		return
//	}
//
//	req.Content = service.EmailServicer.VerifyEmailByCode(ctx, req.Code)
//	if len(req.Content) == 0 {
//		controller.HandleResponse(ctx, errors.New(constants.EmailVerifyURLExpired), nil)
//		return
//	}
//
//	err := service.UserServicer.UserUnsubscribeNotification(ctx, req)
//	controller.HandleResponse(ctx, err, nil)
//}

func (uc *UserController) SearchUserListByName(ctx *gin.Context) {
	req := &types.GetOtherUserInfoByUsernameReq{}
	if controller.BindAndCheck(ctx, req) {
		return
	}
	resp, err := uc.userService.SearchUserListByName(ctx, req)
	if err != nil {
		controller.HandleResponse(ctx, constants.InternalErrCode, constants.EmailOrPasswordWrong, nil)
		return
	}
	controller.HandleResponse(ctx, constants.SuccessCode, constants.Success, resp)
}
