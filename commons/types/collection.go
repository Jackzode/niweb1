package types

import "time"

// Collection collection
type Collection struct {
	ID                    string    `xorm:"not null pk default 0 BIGINT(20) id"`
	CreatedAt             time.Time `xorm:"created not null default CURRENT_TIMESTAMP TIMESTAMP created_at"`
	UpdatedAt             time.Time `xorm:"updated not null default CURRENT_TIMESTAMP TIMESTAMP updated_at"`
	UserID                string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	ObjectID              string    `xorm:"not null default 0 BIGINT(20) object_id"`
	UserCollectionGroupID string    `xorm:"not null default 0 BIGINT(20) user_collection_group_id"`
}

type CollectionSearch struct {
	Page     int `json:"page" form:"page"`           //Query number of pages
	PageSize int `json:"page_size" form:"page_size"` //Search page size
	UserID   string
}

// TableName collection table name
func (Collection) TableName() string {
	return "collection"
}

// CollectionGroup collection group
type CollectionGroup struct {
	ID           string    `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt    time.Time `xorm:"created not null default CURRENT_TIMESTAMP TIMESTAMP created_at"`
	UpdatedAt    time.Time `xorm:"updated not null default CURRENT_TIMESTAMP TIMESTAMP updated_at"`
	UserID       string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	Name         string    `xorm:"not null default '' VARCHAR(50) name"`
	DefaultGroup int       `xorm:"not null default 1 INT(11) default_group"`
}

// TableName collection group table name
func (CollectionGroup) TableName() string {
	return "collection_group"
}

// CollectionSwitchReq switch collection request
type CollectionSwitchReq struct {
	ObjectID string `validate:"required" json:"object_id" form:"object_id"`
	GroupID  string `validate:"omitempty" json:"group_id" form:"group_id"`
	Bookmark bool   `validate:"omitempty" json:"bookmark" form:"bookmark"`
	UserID   string `json:"-"`
}

// CollectionSwitchResp switch collection response
type CollectionSwitchResp struct {
	ObjectCollectionCount int64 `json:"object_collection_count"`
}

// AddCollectionGroupReq add collection group request
type AddCollectionGroupReq struct {
	//
	UserID int64 `validate:"required" comment:"" json:"user_id"`
	// the collection group name
	Name string `validate:"required,gt=0,lte=50" comment:"the collection group name" json:"name"`
	// mark this group is default, default 1
	DefaultGroup int `validate:"required" comment:"mark this group is default, default 1" json:"default_group"`
	//
	CreateTime time.Time `validate:"required" comment:"" json:"create_time"`
	//
	UpdateTime time.Time `validate:"required" comment:"" json:"update_time"`
}

// UpdateCollectionGroupReq update collection group request
type UpdateCollectionGroupReq struct {
	//
	ID int64 `validate:"required" comment:"" json:"id"`
	//
	UserID int64 `validate:"omitempty" comment:"" json:"user_id"`
	// the collection group name
	Name string `validate:"omitempty,gt=0,lte=50" comment:"the collection group name" json:"name"`
	// mark this group is default, default 1
	DefaultGroup int `validate:"omitempty" comment:"mark this group is default, default 1" json:"default_group"`
	//
	CreateTime time.Time `validate:"omitempty" comment:"" json:"create_time"`
	//
	UpdateTime time.Time `validate:"omitempty" comment:"" json:"update_time"`
}

// GetCollectionGroupResp get collection group response
type GetCollectionGroupResp struct {
	//
	ID int64 `json:"id"`
	//
	UserID int64 `json:"user_id"`
	// the collection group name
	Name string `json:"name"`
	// mark this group is default, default 1
	DefaultGroup int `json:"default_group"`
	//
	CreateTime time.Time `json:"create_time"`
	//
	UpdateTime time.Time `json:"update_time"`
}
