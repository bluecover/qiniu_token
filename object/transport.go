package object

// Transport for HTTP

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bluecover/qiniu_token/base"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	endpoints := MakeServerEndpoints(s, logger)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getUploadTokenHandler := kithttp.NewServer(
		endpoints.GetUploadTokenEndpoint,
		decodeGetUploadTokenRequest,
		encodeResponse,
		options...,
	)
	getAccessSecretsHandler := kithttp.NewServer(
		endpoints.GetAccessSecretsEndpoint,
		decodeGetAccessSecretsRequest,
		encodeResponse,
		options...,
	)
	getPrivateURLHandler := kithttp.NewServer(
		endpoints.GetPrivateURLEndpoint,
		decodeGetPrivateURLRequest,
		encodeResponse,
		options...,
	)
	getObjectHandler := kithttp.NewServer(
		endpoints.GetObjectEndpoint,
		decodeGetObjectRequest,
		encodeResponse,
		options...,
	)
	getAllObjectsHandler := kithttp.NewServer(
		endpoints.GetAllObjectsEndpoint,
		decodeGetAllObjectsRequest,
		encodeResponse,
		options...,
	)
	addObjectReferenceHandler := kithttp.NewServer(
		endpoints.AddObjectReferenceEndPoint,
		decodeAddObjectReferenceRequest,
		encodeResponse,
		options...,
	)
	removeObjectReferenceHandler := kithttp.NewServer(
		endpoints.RemoveObjectReferenceEndPoint,
		decodeRemoveObjectReferenceRequest,
		encodeResponse,
		options...,
	)
	deleteObjectHandler := kithttp.NewServer(
		endpoints.DeleteObjectEndpoint,
		decodeDeleteObjectRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/oss/upload/token", getUploadTokenHandler).Methods("GET").Queries("cloud", "{cloud}", "category", "{category}", "user", "{user}")
	r.Handle("/v1/oss/secrets", getAccessSecretsHandler).Methods("GET").Queries("cloud", "{cloud}", "bucket", "{bucket}", "options", "{options}")
	r.Handle("/v1/oss/download/url", getPrivateURLHandler).Methods("GET").Queries("cloud", "{cloud}", "domain", "{domain}", "key", "{key}")
	r.Handle("/v1/oss/get/{id:[0-9]+}", getObjectHandler).Methods("GET")
	r.Handle("/v1/oss/all", getAllObjectsHandler).Methods("GET")
	r.Handle("/v1/oss/addref", addObjectReferenceHandler).Methods("POST")
	r.Handle("/v1/oss/delref", removeObjectReferenceHandler).Methods("POST")
	r.Handle("/v1/oss/del", deleteObjectHandler).Methods("POST")

	return r
}

func decodeGetUploadTokenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	cloud := vars["cloud"]
	if len(cloud) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty cloud"))
	}
	category := vars["category"]
	if len(category) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty category"))
	}
	user := vars["user"]
	if len(user) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty user"))
	}
	return getUploadTokenRequest{
		Cloud:    cloud,
		Category: category,
		User:     user,
	}, nil
}

func decodeGetAccessSecretsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	cloud := vars["cloud"]
	if len(cloud) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty cloud"))
	}
	bucket := vars["bucket"]
	if len(bucket) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty bucket"))
	}
	options := vars["options"]
	return getAccessSecretsRequest{
		Cloud:   cloud,
		Bucket:  bucket,
		Options: options,
	}, nil
}

func decodeGetPrivateURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	cloud := vars["cloud"]
	if len(cloud) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty cloud"))
	}
	domain := vars["domain"]
	if len(domain) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty domain"))
	}
	key := vars["key"]
	if len(key) == 0 {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("empty key"))
	}
	return getPrivateURLRequest{
		Cloud:  cloud,
		Domain: domain,
		Key:    key,
	}, nil
}

func decodeAddObjectReferenceRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeAddObjectReferenceRequest"))
	}
	var req addObjectReferenceRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeAddObjectReferenceRequest"))
	}
	return req, nil
}

func decodeRemoveObjectReferenceRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeRemoveObjectReferenceRequest"))
	}
	var req removeObjectReferenceRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeRemoveObjectReferenceRequest"))
	}
	return req, nil
}

func decodeDeleteObjectRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeDeleteObjectRequest"))
	}
	var req deleteObjectRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, base.NewAppError(ErrInvalidBody, errors.Wrap(err, "decodeDeleteObjectRequest"))
	}
	return req, nil
}

func decodeGetObjectRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, base.NewAppError(ErrMissingParameter, fmt.Errorf("missing id"))
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, base.NewAppError(ErrInvalidParameter, fmt.Errorf("invalid id"))
	}
	return getObjectRequest{
		ID: uint(id),
	}, nil
}

func decodeGetAllObjectsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	var status base.Status
	appErr, ok := err.(*base.AppError)
	if ok {
		status.Code = appErr.Code
		status.Msg = appErr.Err.Error()
	} else {
		status.Code = ErrUnknown
		status.Msg = err.Error()
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": status,
	})
}
