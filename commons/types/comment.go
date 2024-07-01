package types

import (
	"database/sql"
	"fmt"
	"github.com/Jackzode/painting/commons/utils"
	"github.com/jinzhu/copier"
	"time"
)

// Comment comment
type Comment struct {
	ID             string        `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt      time.Time     `xorm:"created TIMESTAMP created_at"`
	UpdatedAt      time.Time     `xorm:"updated TIMESTAMP updated_at"`
	UserID         string        `xorm:"not null default 0 BIGINT(20) user_id"`
	ReplyUserID    sql.NullInt64 `xorm:"BIGINT(20) reply_user_id"`
	ReplyCommentID sql.NullInt64 `xorm:"BIGINT(20) reply_comment_id"`
	ObjectID       string        `xorm:"not null default 0 BIGINT(20) INDEX object_id"`
	QuestionID     string        `xorm:"not null default 0 BIGINT(20) question_id"`
	VoteCount      int           `xorm:"not null default 0 INT(11) vote_count"`
	Status         int           `xorm:"not null default 0 TINYINT(4) status"`
	OriginalText   string        `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText     string        `xorm:"not null MEDIUMTEXT parsed_text"`
}

// TableName comment table name
func (c *Comment) TableName() string {
	return "comment"
}

// GetReplyUserID get reply user id
func (c *Comment) GetReplyUserID() string {
	if c.ReplyUserID.Valid {
		return fmt.Sprintf("%d", c.ReplyUserID.Int64)
	}
	return ""
}

// GetReplyCommentID get reply comment id
func (c *Comment) GetReplyCommentID() string {
	if c.ReplyCommentID.Valid {
		return fmt.Sprintf("%d", c.ReplyCommentID.Int64)
	}
	return ""
}

// SetReplyUserID set reply user id
func (c *Comment) SetReplyUserID(str string) {
	if len(str) > 0 {
		c.ReplyUserID = sql.NullInt64{Int64: utils.Str2Int64(str), Valid: true}
	} else {
		c.ReplyUserID = sql.NullInt64{Valid: false}
	}
}

// SetReplyCommentID set reply comment id
func (c *Comment) SetReplyCommentID(str string) {
	if len(str) > 0 {
		c.ReplyCommentID = sql.NullInt64{Int64: utils.Str2Int64(str), Valid: true}
	} else {
		c.ReplyCommentID = sql.NullInt64{Valid: false}
	}
}

type CommentQuery struct {
	// object id
	ObjectID string
	// query condition
	QueryCond string
	// user id
	UserID   string
	Page     int
	PageSize int
}

// AddCommentReq add comment request
type AddCommentReq struct {
	// object id
	ObjectID string `validate:"required" json:"object_id"`
	// reply comment id
	ReplyCommentID string `validate:"omitempty" json:"reply_comment_id"`
	// original comment content
	OriginalText string `validate:"required,notblank,gte=2,lte=600" json:"original_text"`
	// parsed comment content
	ParsedText string `json:"-"`
	// @ user id list
	MentionUsernameList []string `validate:"omitempty" json:"mention_username_list"`
	// user id
	UserID string `json:"-"`
	// whether user can add it
	CanAdd bool `json:"-"`
	// whether user can edit it
	CanEdit bool `json:"-"`
	// whether user can delete it
	CanDelete   bool   `json:"-"`
	CaptchaID   string `json:"captcha_id"` // captcha_id
	CaptchaCode string `json:"captcha_code"`
}

// RemoveCommentReq remove comment
type RemoveCommentReq struct {
	// comment id
	CommentID string `validate:"required" json:"comment_id"`
	// user id
	UserID      string `json:"-"`
	CaptchaID   string `json:"captcha_id"` // captcha_id
	CaptchaCode string `json:"captcha_code"`
}

// UpdateCommentReq update comment request
type UpdateCommentReq struct {
	// comment id
	CommentID string `validate:"required" json:"comment_id"`
	// original comment content
	OriginalText string `validate:"required,notblank,gte=2,lte=600" json:"original_text"`
	// parsed comment content
	ParsedText string `json:"-"`
	// user id
	UserID  string `json:"-"`
	IsAdmin bool   `json:"-"`

	// whether user can edit it
	CanEdit bool `json:"-"`

	// whether user can delete it
	CaptchaID   string `json:"captcha_id"` // captcha_id
	CaptchaCode string `json:"captcha_code"`
}

type UpdateCommentResp struct {
	// comment id
	CommentID string `json:"comment_id"`
	// original comment content
	OriginalText string `json:"original_text"`
	// parsed comment content
	ParsedText string `json:"parsed_text"`
}

// GetCommentListReq get comment list all request
type GetCommentListReq struct {
	// user id
	UserID int64 `validate:"omitempty" comment:"user id" form:"user_id"`
	// reply user id
	ReplyUserID int64 `validate:"omitempty" comment:"reply user id" form:"reply_user_id"`
	// reply comment id
	ReplyCommentID int64 `validate:"omitempty" comment:"reply comment id" form:"reply_comment_id"`
	// object id
	ObjectID int64 `validate:"omitempty" comment:"object id" form:"object_id"`
	// user vote amount
	VoteCount int `validate:"omitempty" comment:"user vote amount" form:"vote_count"`
	// comment status(available: 0; deleted: 10)
	Status int `validate:"omitempty" comment:"comment status(available: 0; deleted: 10)" form:"status"`
	// original comment content
	OriginalText string `validate:"omitempty" comment:"original comment content" form:"original_text"`
	// parsed comment content
	ParsedText string `validate:"omitempty" comment:"parsed comment content" form:"parsed_text"`
}

// GetCommentWithPageReq get comment list page request
type GetCommentWithPageReq struct {
	// page
	Page int `validate:"omitempty,min=1" form:"page"`
	// page size
	PageSize int `validate:"omitempty,min=1" form:"page_size"`
	// object id
	ObjectID string `validate:"required" form:"object_id"`
	// comment id
	CommentID string `validate:"omitempty" form:"comment_id"`
	// query condition
	QueryCond string `validate:"omitempty,oneof=vote" form:"query_cond"`
	// user id
	UserID string `json:"-"`
	// whether user can edit it
	CanEdit bool `json:"-"`
	// whether user can delete it
	CanDelete bool `json:"-"`
}

// GetCommentReq get comment list page request
type GetCommentReq struct {
	// object id
	ID string `validate:"required" form:"id"`
	// user id
	UserID string `json:"-"`
	// whether user can edit it
	CanEdit bool `json:"-"`
	// whether user can delete it
	CanDelete bool `json:"-"`
}

// GetCommentResp comment response
type GetCommentResp struct {
	// comment id
	CommentID string `json:"comment_id"`
	// create time
	CreatedAt int64 `json:"created_at"`

	// object id
	ObjectID string `json:"object_id"`
	// user vote amount
	VoteCount int `json:"vote_count"`
	// current user if already vote this comment
	IsVote bool `json:"is_vote"`
	// original comment content
	OriginalText string `json:"original_text"`
	// parsed comment content
	ParsedText string `json:"parsed_text"`

	// user id
	UserID string `json:"user_id"`
	// username
	Username string `json:"username"`
	// user display name
	UserDisplayName string `json:"user_display_name"`
	// user avatar
	UserAvatar string `json:"user_avatar"`
	// user status
	UserStatus string `json:"user_status"`

	// reply user id
	ReplyUserID string `json:"reply_user_id"`
	// reply user username
	ReplyUsername string `json:"reply_username"`
	// reply user display name
	ReplyUserDisplayName string `json:"reply_user_display_name"`
	// reply comment id
	ReplyCommentID string `json:"reply_comment_id"`
	// reply user status
	ReplyUserStatus string `json:"reply_user_status"`

	// MemberActions
	//MemberActions []*PermissionMemberAction `json:"member_actions"`
}

func (r *GetCommentResp) SetFromComment(comment *Comment) {
	_ = copier.Copy(r, comment)
	r.CommentID = comment.ID
	r.CreatedAt = comment.CreatedAt.Unix()
	r.ReplyUserID = comment.GetReplyUserID()
	r.ReplyCommentID = comment.GetReplyCommentID()
}

// GetCommentPersonalWithPageReq get comment list page request
type GetCommentPersonalWithPageReq struct {
	// page
	Page int `validate:"omitempty,min=1" form:"page"`
	// page size
	PageSize int `validate:"omitempty,min=1" form:"page_size"`
	// username
	Username string `validate:"omitempty,gt=0,lte=100" form:"username"`
	// user id
	UserID string `json:"-"`
}

// GetCommentPersonalWithPageResp comment response
type GetCommentPersonalWithPageResp struct {
	// comment id
	CommentID string `json:"comment_id"`
	// create time
	CreatedAt int64 `json:"created_at"`
	// object id
	ObjectID string `json:"object_id"`
	// question id
	QuestionID string `json:"question_id"`
	// answer id
	AnswerID string `json:"answer_id"`
	// object type
	ObjectType string `json:"object_type" enums:"question,answer,tag,comment"`
	// title
	Title string `json:"title"`
	// url title
	UrlTitle string `json:"url_title"`
	// content
	Content string `json:"content"`
}
