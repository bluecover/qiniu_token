package object

import (
	"context"

	"github.com/bluecover/qiniu_token/base"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
)

// Endpoints collects all of the endpoints that compose a Object service.
type Endpoints struct {
	GetUploadTokenEndpoint        endpoint.Endpoint
	GetPrivateURLEndpoint         endpoint.Endpoint
	GetAccessSecretsEndpoint      endpoint.Endpoint
	AddObjectReferenceEndPoint    endpoint.Endpoint
	RemoveObjectReferenceEndPoint endpoint.Endpoint
	GetObjectEndpoint             endpoint.Endpoint
	GetAllObjectsEndpoint         endpoint.Endpoint
	DeleteObjectEndpoint          endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeServerEndpoints(s Service, logger log.Logger) Endpoints {
	var (
		debug = viper.GetBool("debug")
	)
	if debug {
		return Endpoints{
			GetUploadTokenEndpoint:        LoggingMiddleware(log.With(logger, "method", "GetUploadToken"))(MakeGetUploadTokenEndpoint(s)),
			GetAccessSecretsEndpoint:      LoggingMiddleware(log.With(logger, "method", "GetAccessSecrets"))(MakeGetAccessSecretsEndpoint(s)),
			GetPrivateURLEndpoint:         LoggingMiddleware(log.With(logger, "method", "GetPrivateURL"))(MakeGetPrivateURLEndpoint(s)),
			AddObjectReferenceEndPoint:    LoggingMiddleware(log.With(logger, "method", "AddObjectReference"))(MakeAddObjectReferenceEndPoint(s)),
			RemoveObjectReferenceEndPoint: LoggingMiddleware(log.With(logger, "method", "RemoveObjectReference"))(MakeRemoveObjectReferenceEndPoint(s)),
			GetObjectEndpoint:             LoggingMiddleware(log.With(logger, "method", "GetObjec"))(MakeGetObjectEndpoint(s)),
			GetAllObjectsEndpoint:         LoggingMiddleware(log.With(logger, "method", "GetAllObjects"))(MakeGetAllObjectsEndpoint(s)),
			DeleteObjectEndpoint:          LoggingMiddleware(log.With(logger, "method", "DeleteObject"))(MakeDeleteObjectEndpoint(s)),
		}
	}
	return Endpoints{
		GetUploadTokenEndpoint:        MakeGetUploadTokenEndpoint(s),
		GetAccessSecretsEndpoint:      MakeGetAccessSecretsEndpoint(s),
		GetPrivateURLEndpoint:         MakeGetPrivateURLEndpoint(s),
		AddObjectReferenceEndPoint:    MakeAddObjectReferenceEndPoint(s),
		RemoveObjectReferenceEndPoint: MakeRemoveObjectReferenceEndPoint(s),
		GetObjectEndpoint:             MakeGetObjectEndpoint(s),
		GetAllObjectsEndpoint:         MakeGetAllObjectsEndpoint(s),
		DeleteObjectEndpoint:          MakeDeleteObjectEndpoint(s),
	}
}

type getUploadTokenRequest struct {
	Cloud    string `json:"cloud"`
	Category string `json:"category"`
	User     string `json:"user"`
}

type getUploadTokenResponse struct {
	Data   UploadToken    `json:"data"`
	Status base.Status    `json:"status"`
	Err    *base.AppError `json:"-"`
}

func (r getUploadTokenResponse) error() *base.AppError { return r.Err }

// MakeGetUploadTokenEndpoint returns an endpoint via the passed service.
func MakeGetUploadTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUploadTokenRequest)
		token, err := s.GetUploadToken(ctx, req.Cloud, req.Category, req.User)
		return getUploadTokenResponse{
			Data:   token,
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type getPrivateURLRequest struct {
	Cloud  string `json:"cloud"`
	Domain string `json:"domain"`
	Key    string `json:"key"`
}

type getPrivateURLResponseData struct {
	URL PrivateURL `json:"url"`
}

type getPrivateURLResponse struct {
	Data   getPrivateURLResponseData `json:"data"`
	Status base.Status               `json:"status"`
	Err    *base.AppError            `json:"-"`
}

func (r getPrivateURLResponse) error() *base.AppError { return r.Err }

// MakeGetPrivateURLEndpoint returns an endpoint via the passed service.
func MakeGetPrivateURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPrivateURLRequest)
		url, err := s.GetPrivateURL(ctx, req.Cloud, req.Domain, req.Key)
		return getPrivateURLResponse{
			Data:   getPrivateURLResponseData{URL: url},
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type getAccessSecretsRequest struct {
	Cloud   string `json:"cloud"`
	Bucket  string `json:"bucket"`
	Options string `json:"options"`
}

type getAccessSecretsResponseData struct {
	Secrets AccessSecrets `json:"secrets"`
}

type getAccessSecretsResponse struct {
	Data   getAccessSecretsResponseData `json:"data"`
	Status base.Status                  `json:"status"`
	Err    *base.AppError               `json:"-"`
}

func (r getAccessSecretsResponse) error() *base.AppError { return r.Err }

// MakeGetAccessSecretsEndpoint returns an endpoint via the passed service.
func MakeGetAccessSecretsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccessSecretsRequest)
		secrets, err := s.GetAccessSecrets(ctx, req.Cloud, req.Bucket, req.Options)
		return getAccessSecretsResponse{
			Data:   getAccessSecretsResponseData{Secrets: secrets},
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type addObjectReferenceRequest struct {
	UserID uint       `json:"userID"`
	Tag    string     `json:"tag"`
	Object ObjectInfo `json:"object"`
}

type addObjectReferenceResponse struct {
	Status base.Status    `json:"status"`
	Err    *base.AppError `json:"-"`
}

func (r addObjectReferenceResponse) error() *base.AppError { return r.Err }

// MakeAddObjectReferenceEndPoint returns an endpoint via the passed service.
func MakeAddObjectReferenceEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addObjectReferenceRequest)
		err := s.AddObjectReference(ctx, req.UserID, req.Tag, req.Object)
		return addObjectReferenceResponse{
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type removeObjectReferenceRequest struct {
	UserID   uint   `json:"userID"`
	ObjectID uint   `json:"objectID"`
	Tag      string `json:"tag"`
}

type removeObjectReferenceResponse struct {
	Status base.Status    `json:"status"`
	Err    *base.AppError `json:"-"`
}

func (r removeObjectReferenceResponse) error() *base.AppError { return r.Err }

// MakeRemoveObjectReferenceEndPoint returns an endpoint via the passed service.
func MakeRemoveObjectReferenceEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeObjectReferenceRequest)
		err := s.RemoveObjectReference(ctx, req.UserID, req.ObjectID, req.Tag)
		return removeObjectReferenceResponse{
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type getObjectRequest struct {
	ID uint `json:"id"`
}

type getObjectResponseData struct {
	Object ObjectInfo `json:"object"`
}

type getObjectResponse struct {
	Data   getObjectResponseData `json:"data"`
	Status base.Status           `json:"status"`
	Err    *base.AppError        `json:"-"`
}

func (r getObjectResponse) error() *base.AppError { return r.Err }

// MakeGetObjectEndpoint returns an endpoint via the passed service.
func MakeGetObjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getObjectRequest)
		object, err := s.GetObject(ctx, req.ID)
		return getObjectResponse{
			Data:   getObjectResponseData{Object: object},
			Status: base.SuccessStatus,
			Err:    err,
		}, nil
	}
}

type getAllObjectResponseData struct {
	Objects []ObjectInfo `json:"objects"`
}

type getAllObjectsResponse struct {
	Data   getAllObjectResponseData `json:"data"`
	Status base.Status              `json:"status"`
	Err    *base.AppError           `json:"-"`
}

func (r getAllObjectsResponse) error() *base.AppError { return r.Err }

// MakeGetAllObjectsEndpoint returns an endpoint via the passed service.
func MakeGetAllObjectsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		objects, err := s.GetAllObjects(ctx)
		return getAllObjectsResponse{
			Data:   getAllObjectResponseData{Objects: objects},
			Status: base.SuccessStatus, Err: err,
		}, nil
	}
}

type deleteObjectRequest struct {
	ID uint `json:"id"`
}

type deleteObjectResponse struct {
	Status base.Status    `json:"status"`
	Err    *base.AppError `json:"-"`
}

func (r deleteObjectResponse) error() *base.AppError { return r.Err }

// MakeDeleteObjectEndpoint returns an endpoint via the passed service.
func MakeDeleteObjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteObjectRequest)
		e := s.DeleteObject(ctx, req.ID)
		return deleteObjectResponse{Status: base.SuccessStatus, Err: e}, nil
	}
}
