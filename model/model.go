package model

// Model definitions using gorm

import (
	"time"
)

// Object represents Object model.
type Object struct {
	ID          uint      `gorm:"column:id;primary_key;auto_increment"`
	Cloud       string    `gorm:"column:cloud;type:varchar(8);not null;unique_index:objcet_cloud_bucket_key_unique"`
	Bucket      string    `gorm:"column:bucket;type:varchar(64);not null;unique_index:objcet_cloud_bucket_key_unique"`
	Key         string    `gorm:"column:key;type:varchar(128);not null;unique_index:objcet_cloud_bucket_key_unique"`
	Etag        string    `gorm:"column:etag;type:varchar(32)"`
	MimeType    string    `gorm:"column:mime_type;type:varchar(16)"`
	Size        uint      `gorm:"column:size"`
	Status      int       `gorm:"column:status;type:tinyint(1);not null"`
	CreatedTime time.Time `gorm:"column:created_time;type:timestamp"`
}

// TableName defines table name in database.
func (Object) TableName() string {
	return "oss"
}

// ObjectRef represents ObjectRef model.
type ObjectRef struct {
	ID          uint      `gorm:"column:id;primary_key;auto_increment"`
	UserID      uint      `gorm:"column:user_id;not null"`
	ObjectID    uint      `gorm:"column:object_id;not null"`
	Tag         string    `gorm:"column:tag;type:varchar(32);not null"`
	Status      int       `gorm:"column:status;type:tinyint(1)"`
	CreatedTime time.Time `gorm:"column:created_time;type:timestamp"`
}

// TableName defines table name in database.
func (ObjectRef) TableName() string {
	return "oss_ref"
}
