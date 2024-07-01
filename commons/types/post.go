package types

import (
	"time"
)

// RemoveQuestionReq delete question request
type RemoveQuestionReq struct {
	// question id
	ID          string `validate:"required" json:"id"`
	UserID      string `json:"-" ` // user_id
	IsAdmin     bool   `json:"-"`
	CaptchaID   string `json:"captcha_id"` // captcha_id
	CaptchaCode string `json:"captcha_code"`
}

type CloseQuestionReq struct {
	ID        string `validate:"required" json:"id"`
	CloseType int    `json:"close_type"` // close_type
	CloseMsg  string `json:"close_msg"`  // close_type
	UserID    string `json:"-"`          // user_id
}

type OperationQuestionReq struct {
	ID        string `validate:"required" json:"id"`
	Operation string `json:"operation"` // operation [pin unpin hide show]
	UserID    string `json:"-"`         // user_id
	CanPin    bool   `json:"-"`
	CanList   bool   `json:"-"`
}

type CloseQuestionMeta struct {
	CloseType int    `json:"close_type"`
	CloseMsg  string `json:"close_msg"`
}

// ReopenQuestionReq reopen question request
type ReopenQuestionReq struct {
	QuestionID string `json:"question_id"`
	UserID     string `json:"-"`
}

type QuestionAdd struct {
	// question title
	Title string `validate:"required,notblank,gte=6,lte=150" json:"title"`
	// content
	Content string `validate:"required,notblank,gte=6,lte=65535" json:"content"`
	// html
	HTML string `json:"-"`
	// user id
	UserID       string `json:"-"`
	CopyRight    string `validate:"required" json:"copyright"`
	AllowReprint string `validate:"required " json:"allow_reprint"`
	AllowComment string `validate:"required" json:"allow_comment"`
	Feeds        string `validate:"required " json:"feeds"`
	CaptchaID    string `json:"captcha_id"` // captcha_id
	CaptchaCode  string `json:"captcha_code"`
}

type QuestionAddByAnswer struct {
	// question title
	Title string `validate:"required,notblank,gte=6,lte=150" json:"title"`
	// content
	Content string `validate:"required,notblank,gte=6,lte=65535" json:"content"`
	// html
	HTML          string `json:"-"`
	AnswerContent string `validate:"required,notblank,gte=6,lte=65535" json:"answer_content"`
	AnswerHTML    string `json:"-"`
	// tags
	//Tags []*TagItem `validate:"required,dive" json:"tags"`
	// user id
	UserID              string   `json:"-"`
	MentionUsernameList []string `validate:"omitempty" json:"mention_username_list"`
	CaptchaID           string   `json:"captcha_id"` // captcha_id
	CaptchaCode         string   `json:"captcha_code"`
}

type CheckCanQuestionUpdate struct {
	// question id
	ID string `validate:"required" form:"id"`
	// user id
	UserID  string `json:"-"`
	IsAdmin bool   `json:"-"`
}

type QuestionUpdate struct {
	// question id
	ID string `validate:"required" json:"id"`
	// question title
	Title string `validate:"required,notblank,gte=6,lte=150" json:"title"`
	// content
	Content string `validate:"required,notblank,gte=6,lte=65535" json:"content"`
	// html
	HTML       string   `json:"-"`
	InviteUser []string `validate:"omitempty"  json:"invite_user"`
	// tags
	//Tags []*TagItem `validate:"required,dive" json:"tags"`
	// edit summary
	EditSummary string `validate:"omitempty" json:"edit_summary"`
	// user id
	UserID       string `json:"-"`
	NoNeedReview bool   `json:"-"`
	CaptchaID    string `json:"captcha_id"` // captcha_id
	CaptchaCode  string `json:"captcha_code"`
}

type QuestionRecoverReq struct {
	QuestionID string `validate:"required" json:"question_id"`
	UserID     string `json:"-"`
}

type QuestionUpdateInviteUser struct {
	ID          string   `validate:"required" json:"id"`
	InviteUser  []string `validate:"omitempty"  json:"invite_user"`
	UserID      string   `json:"-"`
	CaptchaID   string   `json:"captcha_id"` // captcha_id
	CaptchaCode string   `json:"captcha_code"`
}

type QuestionBaseInfo struct {
	ID              string `json:"id" `
	Title           string `json:"title"`
	UrlTitle        string `json:"url_title"`
	ViewCount       int    `json:"view_count"`
	AnswerCount     int    `json:"answer_count"`
	CollectionCount int    `json:"collection_count"`
	FollowCount     int    `json:"follow_count"`
	Status          string `json:"status"`
	AcceptedAnswer  bool   `json:"accepted_answer"`
}

type QuestionInfo struct {
	ID          string `json:"id" `
	Title       string `json:"title"`
	UrlTitle    string `json:"url_title"`
	Content     string `json:"content"`
	HTML        string `json:"html"`
	Description string `json:"description"`
	//Tags                 []*TagResp     `json:"tags"`
	ViewCount            int            `json:"view_count"`
	UniqueViewCount      int            `json:"unique_view_count"`
	VoteCount            int            `json:"vote_count"`
	AnswerCount          int            `json:"answer_count"`
	CollectionCount      int            `json:"collection_count"`
	FollowCount          int            `json:"follow_count"`
	AcceptedAnswerID     string         `json:"accepted_answer_id"`
	LastAnswerID         string         `json:"last_answer_id"`
	CreateTime           int64          `json:"create_time"`
	UpdateTime           int64          `json:"-"`
	PostUpdateTime       int64          `json:"update_time"`
	QuestionUpdateTime   int64          `json:"edit_time"`
	Pin                  int            `json:"pin"`
	Show                 int            `json:"show"`
	Status               int            `json:"status"`
	Operation            *Operation     `json:"operation,omitempty"`
	UserID               string         `json:"-"`
	LastEditUserID       string         `json:"-"`
	LastAnsweredUserID   string         `json:"-"`
	UserInfo             *UserBasicInfo `json:"user_info"`
	UpdateUserInfo       *UserBasicInfo `json:"update_user_info,omitempty"`
	LastAnsweredUserInfo *UserBasicInfo `json:"last_answered_user_info,omitempty"`
	Answered             bool           `json:"answered"`
	FirstAnswerId        string         `json:"first_answer_id"`
	Collected            bool           `json:"collected"`
	VoteStatus           string         `json:"vote_status"`
	IsFollowed           bool           `json:"is_followed"`
}

// UpdateQuestionResp update question resp
type UpdateQuestionResp struct {
	UrlTitle      string `json:"url_title"`
	WaitForReview bool   `json:"wait_for_review"`
}

type AdminQuestionInfo struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	VoteCount        int            `json:"vote_count"`
	AnswerCount      int            `json:"answer_count"`
	AcceptedAnswerID string         `json:"accepted_answer_id"`
	CreateTime       int64          `json:"create_time"`
	UpdateTime       int64          `json:"update_time"`
	EditTime         int64          `json:"edit_time"`
	UserID           string         `json:"-" `
	UserInfo         *UserBasicInfo `json:"user_info"`
}

type Operation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Msg         string `json:"msg"`
	Time        int64  `json:"time"`
	Level       string `json:"level"`
}

type GetCloseTypeResp struct {
	// report name
	Name string `json:"name"`
	// report description
	Description string `json:"description"`
	// report source
	Source string `json:"source"`
	// report type
	Type int `json:"type"`
	// is have content
	HaveContent bool `json:"have_content"`
	// content type
	ContentType string `json:"content_type"`
}

type UserAnswerInfo struct {
	AnswerID     string `json:"answer_id"`
	QuestionID   string `json:"question_id"`
	Accepted     int    `json:"accepted"`
	VoteCount    int    `json:"vote_count"`
	CreateTime   int    `json:"create_time"`
	UpdateTime   int    `json:"update_time"`
	QuestionInfo struct {
		Title    string        `json:"title"`
		UrlTitle string        `json:"url_title"`
		Tags     []interface{} `json:"tags"`
	} `json:"question_info"`
}

type UserQuestionInfo struct {
	ID               string        `json:"question_id"`
	Title            string        `json:"title"`
	UrlTitle         string        `json:"url_title"`
	VoteCount        int           `json:"vote_count"`
	Tags             []interface{} `json:"tags"`
	ViewCount        int           `json:"view_count"`
	AnswerCount      int           `json:"answer_count"`
	CollectionCount  int           `json:"collection_count"`
	CreatedAt        int64         `json:"created_at"`
	AcceptedAnswerID string        `json:"accepted_answer_id"`
	Status           string        `json:"status"`
}

// QuestionPageReq query questions page
type QuestionPageReq struct {
	Page      int    `validate:"omitempty,min=1" form:"page"`
	PageSize  int    `validate:"omitempty,min=1" form:"pageSize"`
	OrderCond string `validate:"omitempty,oneof=newest active frequent score unanswered recommend" form:"cond"`
	Tag       string `validate:"omitempty,gt=0,lte=100" form:"tag"`
	Username  string `validate:"omitempty,gt=0,lte=100" form:"username"`
	InDays    int    `validate:"omitempty,min=1" form:"in_days"`

	LoginUserID      string `json:"-"`
	UserIDBeSearched string `json:"-"`
	TagID            string `json:"-"`
}

type QuestionPageResp struct {
	ID          string `json:"id" `
	CreatedAt   int64  `json:"created_at"`
	Title       string `json:"title"`
	UrlTitle    string `json:"url_title"`
	Description string `json:"description"`
	Pin         int    `json:"pin"`  // 1: unpin, 2: pin
	Show        int    `json:"show"` // 0: show, 1: hide
	Status      int    `json:"status"`
	//Tags        []*TagResp `json:"tags"`

	// question statistical information
	ViewCount       int `json:"view_count"`
	UniqueViewCount int `json:"unique_view_count"`
	VoteCount       int `json:"vote_count"`
	AnswerCount     int `json:"answer_count"`
	CollectionCount int `json:"collection_count"`
	FollowCount     int `json:"follow_count"`

	// answer information
	AcceptedAnswerID   string    `json:"accepted_answer_id"`
	LastAnswerID       string    `json:"last_answer_id"`
	LastAnsweredUserID string    `json:"-"`
	LastAnsweredAt     time.Time `json:"-"`
	//
	AuthorID   string     `json:"authorID"`
	AuthorInfo AuthorInfo `json:"authorInfo"`
	//
	Content string `json:"content"`

	// operator information
	//OperatedAt    int64                     `json:"operated_at"`
	//Operator      *QuestionPageRespOperator `json:"operator"`
	//OperationType string                    `json:"operation_type"`

}

type AuthorInfo struct {
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Rank        int    `json:"rank"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type QuestionPageRespOperator struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Rank        int    `json:"rank"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Avatar      string `json:"avatar"`
}

type AdminQuestionPageReq struct {
	Page        int    `validate:"omitempty,min=1" form:"page"`
	PageSize    int    `validate:"omitempty,min=1" form:"page_size"`
	StatusCond  string `validate:"omitempty,oneof=normal closed deleted" form:"status"`
	Query       string `validate:"omitempty,gt=0,lte=100" json:"query" form:"query" `
	Status      int    `json:"-"`
	LoginUserID string `json:"-"`
}

// AdminAnswerPageReq admin answer page req
type AdminAnswerPageReq struct {
	Page          int    `validate:"omitempty,min=1" form:"page"`
	PageSize      int    `validate:"omitempty,min=1" form:"page_size"`
	StatusCond    string `validate:"omitempty,oneof=normal deleted" form:"status"`
	Query         string `validate:"omitempty,gt=0,lte=100" form:"query"`
	QuestionID    string `validate:"omitempty,gt=0,lte=24" form:"question_id"`
	QuestionTitle string `json:"-"`
	AnswerID      string `json:"-"`
	Status        int    `json:"-"`
	LoginUserID   string `json:"-"`
}

type AdminUpdateQuestionStatusReq struct {
	QuestionID string `validate:"required" json:"question_id"`
	Status     string `validate:"required,oneof=available closed deleted" json:"status"`
	UserID     string `json:"-"`
}

type PersonalQuestionPageReq struct {
	Page        int    `validate:"omitempty,min=1" form:"page"`
	PageSize    int    `validate:"omitempty,min=1" form:"pageSize"`
	OrderCond   string `validate:"omitempty,oneof=newest active frequent score unanswered" form:"cond"`
	Username    string `validate:"omitempty,gt=0,lte=100" form:"username"`
	LoginUserID string `json:"-"`
}

type PersonalAnswerPageReq struct {
	Page        int    `validate:"omitempty,min=1" form:"page"`
	PageSize    int    `validate:"omitempty,min=1" form:"page_size"`
	OrderCond   string `validate:"omitempty,oneof=newest active frequent score unanswered" form:"cond"`
	Username    string `validate:"omitempty,gt=0,lte=100" form:"username"`
	LoginUserID string `json:"-"`
}

type PersonalCollectionPageReq struct {
	Page     int    `validate:"omitempty,min=1" form:"page"`
	PageSize int    `validate:"omitempty,min=1" form:"pageSize"`
	UserID   string `json:"-"`
}

type Question struct {
	ID               string    `xorm:"not null pk BIGINT(20) id"`
	CreatedAt        time.Time `xorm:"not null default CURRENT_TIMESTAMP TIMESTAMP created_at"`
	UpdatedAt        time.Time `xorm:"updated_at TIMESTAMP"`
	UserID           string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	InviteUserID     string    `xorm:"TEXT invite_user_id"`
	LastEditUserID   string    `xorm:"not null default 0 BIGINT(20) last_edit_user_id"`
	Title            string    `xorm:"not null default '' VARCHAR(150) title"`
	OriginalText     string    `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText       string    `xorm:"not null MEDIUMTEXT parsed_text"`
	Pin              int       `xorm:"not null default 1 INT(11) pin"`
	Show             int       `xorm:"not null default 1 INT(11) show"`
	Status           int       `xorm:"not null default 1 INT(11) status"`
	ViewCount        int       `xorm:"not null default 0 INT(11) view_count"`
	UniqueViewCount  int       `xorm:"not null default 0 INT(11) unique_view_count"`
	VoteCount        int       `xorm:"not null default 0 INT(11) vote_count"`
	AnswerCount      int       `xorm:"not null default 0 INT(11) answer_count"`
	CollectionCount  int       `xorm:"not null default 0 INT(11) collection_count"`
	FollowCount      int       `xorm:"not null default 0 INT(11) follow_count"`
	AcceptedAnswerID string    `xorm:"not null default 0 BIGINT(20) accepted_answer_id"`
	LastAnswerID     string    `xorm:"not null default 0 BIGINT(20) last_answer_id"`
	PostUpdateTime   time.Time `xorm:"post_update_time TIMESTAMP"`
	RevisionID       string    `xorm:"not null default 0 BIGINT(20) revision_id"`
	CopyRight        int       `xorm:"not null default 0 INT(11) copyright"`
	AllowReprint     int       `xorm:"not null default 0 INT(11) allow_reprint"`
	AllowComment     int       `xorm:"not null default 0 INT(11) allow_comment"`
	Feeds            int       `xorm:"not null default 0 INT(11) feeds"`
}

// TableName question table name
func (Question) TableName() string {
	return "question"
}
