package types

import (
	"time"
)

type CaptchaCodeReq struct {
	Email string `json:"email" form:"email"`
}

type ActionRecordResp struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaImg string `json:"captcha_img"`
	Verify     bool   `json:"verify"`
}

type UserVerifyEmailReq struct {
	// code
	Code string `validate:"required,gt=0,lte=500" form:"code"`
	// content
	Content string `json:"content"`
}

type UserRegisterReq struct {
	Username string `validate:"required,gt=3,lte=30" json:"username"`
	Email    string `validate:"required,email,gt=0,lte=500" json:"email" `
	Password string `validate:"required,gte=8,lte=32" json:"password"`
	//CaptchaID   string `json:"captcha_id" `
	CaptchaCode string `validate:"required " json:"captcha_code"`
	IP          string `json:"-" `
}

type UserChangeEmailSendCodeReq struct {
	CaptchaID   string `validate:"omitempty,gt=0,lte=500" json:"captcha_id"`
	CaptchaCode string `validate:"omitempty,gt=0,lte=500" json:"captcha_code"`
	Email       string `validate:"required,email,gt=0,lte=500" json:"e_mail"`
	Pass        string `validate:"omitempty,gte=8,lte=32" json:"pass"`
	UserID      string `json:"-"`
}

type UpdateInfoRequest struct {
	DisplayName string `validate:"omitempty,gt=0,lte=30" form:"display_name"`
	Description string `validate:"omitempty,gt=0,lte=4096" form:"description"`
	School      string `validate:"omitempty,gt=0,lte=30" form:"school"`
	Position    string `validate:"omitempty,gt=0,lte=30" form:"position"`
	Username    string `validate:"omitempty,gt=3,lte=30" form:"username"`
	Avatar      string `form:"avatar"`
	Company     string `validate:"omitempty,gt=0,lte=4096" form:"Company"`
	Github      string `form:"github"`
	Website     string `validate:"omitempty,gt=0,lte=500" form:"website"`
	CityId      string `validate:"omitempty" form:"cityId"`
	UserID      string `form:"-"`
	FirstName   string `validate:"omitempty,gt=0,lte=100" form:"firstname"`
	LastName    string `validate:"omitempty,gt=0,lte=100" form:"lastname"`
	Birthday    string `validate:"omitempty,gt=0,lte=100" form:"birthday"`
}

type UserModifyPasswordReq struct {
	OldPass     string `validate:"omitempty,gte=8,lte=32" form:"old_pass"`
	Pass        string `validate:"required,gte=8,lte=32" form:"pass"`
	CaptchaID   string `validate:"omitempty,gt=0,lte=500" form:"captcha_id"`
	CaptchaCode string `validate:"omitempty,gt=0,lte=500" form:"captcha_code"`
	UserID      string `json:"-"`
	AccessToken string `json:"-"`
}

type UserRetrievePassWordRequest struct {
	Email       string `validate:"required,email,gt=0,lte=500" form:"email"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

type UserRePassWordRequest struct {
	Code    string `validate:"required,gt=0,lte=100" form:"code"`
	Pass    string `validate:"required,gt=0,lte=32" form:"pass"`
	Content string `json:"-"`
}

// UserLoginResp get user response
type UserLoginResp struct {
	// user id
	ID string `json:"id"`
	// create time
	CreatedAt int64 `json:"created_at"`
	// last login date
	LastLoginDate int64 `json:"last_login_date"`
	// username
	Username string `json:"username"`
	// email
	EMail string `json:"e_mail"`
	// mail status(1 pass 2 to be verified)
	MailStatus int `json:"mail_status"`
	// notice status(1 on 2off)
	NoticeStatus int `json:"notice_status"`
	// follow count
	FollowCount int `json:"follow_count"`
	// question count
	QuestionCount int `json:"question_count"`
	// display name
	DisplayName string `json:"display_name"`
	// avatar
	Avatar string `json:"avatar"`
	// mobile
	Mobile string `json:"mobile"`
	// bio markdown
	Description string `json:"description"`
	// bio html
	Company string `json:"company"`
	// website
	Website string `json:"website"`
	// location
	CityId string `json:"city_id"`
	// language
	Language string `json:"language"`
	// access token
	Token string `json:"token"`
	// role id
	RoleID int `json:"role_id"`
	// user status
	Status string `json:"status"`
	// user have password
	HavePassword bool `json:"have_password"`
	//school
	School string `json:"school"`
	//github
	Github string `json:"github"`
	//birthday
	Birthday  int64  `json:"birthday"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Position  string `json:"position"`
}

type UserEmailLoginReq struct {
	Email       string `validate:"required,email,gt=0,lte=500" json:"email"`
	Password    string `validate:"required,gte=8,lte=32" json:"password"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

type GetOtherUserInfoByUsernameReq struct {
	Username string `validate:"required,gt=0,lte=500" form:"username"`
	UserID   string `json:"-"`
}

// TableName user table name
func (User) TableName() string {
	return "user"
}

// User user
type User struct {
	ID             string    `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt      time.Time `xorm:"created TIMESTAMP created_at"`
	Birthday       time.Time `xorm:"TIMESTAMP birthday"`
	UpdatedAt      time.Time `xorm:"updated TIMESTAMP updated_at"`
	SuspendedAt    time.Time `xorm:"TIMESTAMP suspended_at"`
	DeletedAt      time.Time `xorm:"TIMESTAMP deleted_at"`
	LastLoginDate  time.Time `xorm:"TIMESTAMP last_login_date"`
	Username       string    `xorm:"not null default '' VARCHAR(50) UNIQUE username"`
	Pass           string    `xorm:"not null default '' VARCHAR(255) pass"`
	EMail          string    `xorm:"not null VARCHAR(100) e_mail"`
	MailStatus     int       `xorm:"not null default 2 TINYINT(4) mail_status"`
	NoticeStatus   int       `xorm:"not null default 2 INT(11) notice_status"`
	FollowCount    int       `xorm:"not null default 0 INT(11) follow_count"`
	AnswerCount    int       `xorm:"not null default 0 INT(11) answer_count"`
	QuestionCount  int       `xorm:"not null default 0 INT(11) question_count"`
	Rank           int       `xorm:"not null default 0 INT(11) rank"`
	Status         int       `xorm:"not null default 1 INT(11) status"`
	AuthorityGroup int       `xorm:"not null default 1 INT(11) authority_group"`
	DisplayName    string    `xorm:"not null default '' VARCHAR(30) display_name"`
	Avatar         string    `xorm:"not null default '' VARCHAR(1024) avatar"`
	Mobile         string    `xorm:"not null VARCHAR(20) mobile"`
	Description    string    `xorm:"not null TEXT description"`
	Company        string    `xorm:"not null TEXT company"`
	Website        string    `xorm:"not null default '' VARCHAR(255) website"`
	CityId         string    `xorm:"not null default '' INT(11) city_id"`
	IPInfo         string    `xorm:"not null default '' VARCHAR(255) ip_info"`
	School         string    `xorm:"school"`
	Language       string    `xorm:"not null default '' VARCHAR(100) language"`
	Position       string    `xorm:"not null default '' VARCHAR(100) position"`
	Github         string    `xorm:"not null default '' VARCHAR(100) github"`
	Firstname      string    `xorm:"VARCHAR(100) firstname"`
	Lastname       string    `xorm:"VARCHAR(100) lastname"`
}

type GetOtherUserInfoByUsername struct {
	// user id
	ID string `json:"id"`
	// create time
	CreatedAt int64 `json:"created_at"`
	// last login date
	LastLoginDate int64 `json:"last_login_date"`
	// username
	Username string `json:"username"`
	// email
	// follow count
	FollowCount int `json:"follow_count"`
	// answer count
	AnswerCount int `json:"answer_count"`
	// question count
	QuestionCount int `json:"question_count"`
	// rank
	Rank int `json:"rank"`
	// display name
	DisplayName string `json:"display_name"`
	// avatar
	Avatar string `json:"avatar"`
	// mobile
	Mobile string `json:"mobile"`
	// bio markdown
	Bio string `json:"bio"`
	// bio html
	BioHTML string `json:"bio_html"`
	// website
	Website string `json:"website"`
	// location
	CityId    string `json:"city_id"`
	Status    string `json:"status"`
	StatusMsg string `json:"status_msg,omitempty"`
}

type UserBasicInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Rank        int    `json:"rank"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
	Website     string `json:"website"`
	CityId      string `json:"city_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
