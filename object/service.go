package object

import (
	"context"
	"time"

	"github.com/bluecover/qiniu_token/base"
)

// Service interface for service.
type Service interface {
	GetUploadToken(ctx context.Context, cloud string, category string, user string) (UploadToken, *base.AppError)
	GetAccessSecrets(ctx context.Context, cloud string, bucket string, options string) (AccessSecrets, *base.AppError)
	GetPrivateURL(ctx context.Context, cloud string, domain string, key string) (PrivateURL, *base.AppError)
	AddObjectReference(ctx context.Context, userID uint, tag string, objInfo ObjectInfo) *base.AppError
	RemoveObjectReference(ctx context.Context, userID uint, objectID uint, tag string) *base.AppError
	GetObject(ctx context.Context, id uint) (ObjectInfo, *base.AppError)
	GetAllObjects(ctx context.Context) ([]ObjectInfo, *base.AppError)
	DeleteObject(ctx context.Context, id uint) *base.AppError
}

// UploadToken represents response data from GetUploadToken
type UploadToken struct {
	Bucket     string    `json:"bucket"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

// AccessSecrets represents response data from GetAccessSecrets
type AccessSecrets struct {
	CloudService    string    `json:"cloudService"`
	AccessKeyID     string    `json:"accessKeyId"`
	AccessKeySecret string    `json:"accessKeySecret"`
	Token           string    `json:"token"`
	Expiration      time.Time `json:"expiration"`
}

// PrivateURL represents response data from GetPrivateURL
type PrivateURL struct {
	URL        string    `json:"url"`
	Expiration time.Time `json:"expiration"`
}

// ObjectInfo represents properties of a object
type ObjectInfo struct {
	Cloud    string `json:"cloud"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Etag     string `json:"etag"`
	MimeType string `json:"mimeType"`
	Size     uint   `json:"size"`
	Status   int    `json:"status"`
}

type errorer interface {
	error() *base.AppError
}

var (
	ErrUnimplemented            = "unimplemented"
	ErrMissingParameter         = "missing parameter"
	ErrInvalidParameter         = "invalid parameter"
	ErrInvalidBody              = "invalid body"
	ErrUnsupportedCloundService = "unsupported cloud service"
	ErrAlreadyExists            = "already exists"
	ErrNotFound                 = "not found"
	ErrUpdateFailed             = "update failed"
	ErrAliyunSTS                = "aliyun STS error"
	ErrModelOperation           = "model operation error"
	ErrUnknown                  = "unknown error"
)
