package types

import "time"

type AddLikeReq struct {
	QuestionID string `validate:"required,gt=0,lte=100" form:"question_id" json:"question_id"`
	Status     bool   `validate:"required" form:"status" json:"status"`
	AuthorID   string `validate:"required,gt=0,lte=100" form:"author_id" json:"author_id"`
	UserId     string
}

type Likes struct {
	ID         string    `xorm:"not null pk BIGINT(20) id"`
	CreatedAt  time.Time `xorm:"not null default CURRENT_TIMESTAMP TIMESTAMP created_at"`
	UpdatedAt  time.Time `xorm:"updated_at TIMESTAMP"`
	QuestionID string    `xorm:"not null BIGINT(20) question_id"`
	UserID     string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	AuthorID   string    `xorm:"not null BIGINT(20) author_id"`
	Status     int       `xorm:"not null default 1 INT(11) status"` //1=like, 0=dislike
}

func (Likes) TableName() string {
	return "likes"
}
