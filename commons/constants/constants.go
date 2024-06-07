package constants

import (
	"time"
)

const (
	StatusErrCode               = 10002
	CaptchaFailedCode           = 10001
	ParamInvalid                = 10000
	AcceptLanguageFlag          = "Accept-Language"
	CaptchaExpiration           = 5 * time.Minute
	TitleRetrievePassWord       = "RetrievePassWord"
	TitleVerifyEmail            = "VerifyEmail "
	TitleChangeEmail            = "ChangeEmail"
	TitleRegisterByEmail        = "RegisterByEmail"
	SSL                         = "SSL"
	CommonRole                  = 1
	UserStatusAvailable         = 1
	UserStatusSuspended         = 9
	UserStatusDeleted           = 10
	DefaultConfigFileName       = "config.yaml"
	ConfigFileDir               = "./conf/"
	EmailEncryption             = "SSL"
	HeaderToken                 = "lawyer-token"
	TokenClaim                  = "TokenClaim"
	TraceID                     = "trace_id"
	Issuer                      = "lawyer-test"
	ExpireDate                  = time.Hour * 24 * 15
	ExpireBuffer          int64 = 1000 * 3600 * 24

	UserNormal              = "normal"
	UserSuspended           = "suspended"
	UserDeleted             = "deleted"
	UserInactive            = "inactive"
	EmailStatusAvailable    = 1
	EmailStatusToBeVerified = 2
	DefaultLanguage         = "en-US"
)
