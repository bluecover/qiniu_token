package callback

// All-in-one for cloud callback handler service

import (
	"context"

	"github.com/bluecover/qiniu_token/base"
	"github.com/go-kit/kit/endpoint"
)

// Service is the interface to handle callbacks from cloud service (OSS, Qiniu, etc).
type Service interface {
	OssPutObjectCallback(ctx context.Context, param OssCallbackParam) (CallbackResult, *base.AppError)
}

// OssCallbackParam represents put-object callback parameters from aliyun OSS.
type OssCallbackParam struct {
	Bucket       string `json:"bucket"`
	Object       string `json:"object"`
	OriginName   string `json:"originName"`
	Etag         string `json:"etag"`
	Size         uint   `json:"size,string"`
	ImageFormat  string `json:"imageInfo.format"`
	ImageWidth   uint   `json:"imageInfo.width,string"`
	ImageHeight  uint   `json:"imageInfo.height,string"`
	AppName      string `json:"appName"`
	AppBusiness  string `json:"appBusiness"`
	AppUserID    uint   `json:"appUserID,string"`
	AppUserToken string `json:"appUserToken"`
}

type CallbackResult struct {
	ObjID       uint   `json:"objID"`
	ObjFilename string `json:"objFilename"`
}

type ossPutObjectCallbackRequest struct {
	Param OssCallbackParam
}

type ossPutObjectCallbackResponse struct {
	Data   CallbackResult `json:"data"`
	Status base.Status    `json:"status"`
	Err    *base.AppError `json:"-"`
}

type errorer interface {
	error() *base.AppError
}

func (r ossPutObjectCallbackResponse) error() *base.AppError {
	return r.Err
}

var (
	ErrModelFunctionFailed    = "model function failed"
	ErrAlreadyExists          = "already exists"
	ErrUserVerificationFailed = "user verification failed"
)

// MakePutObjectEndpoint returns an endpoint via the passed service.
func MakeOssPutObjectCallbackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ossPutObjectCallbackRequest)
		result, err := s.OssPutObjectCallback(ctx, req.Param)
		return ossPutObjectCallbackResponse{Data: result, Status: base.SuccessStatus, Err: err}, nil
	}
}
