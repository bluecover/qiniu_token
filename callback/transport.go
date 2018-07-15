package callback

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bluecover/qiniu_token/base"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func decodeOssPutObjectCallbackRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ossPutObjectCallbackRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeOssPutObjectCallbackRequest:ioutil.ReadAll")
	}
	if err := json.Unmarshal(body, &req.Param); err != nil {
		return nil, errors.Wrap(err, "decodeOssPutObjectCallbackRequest:json.Unmarshal")
	}
	return req, nil
}

func encodeOssPutObjectCallbackResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err *base.AppError, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.Code {
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": base.Status{Code: err.Code, Msg: err.Error()},
	})
}

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoint := MakeOssPutObjectCallbackEndpoint(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	router.Methods("POST").Path("/callback/oss-put-object").Handler(kithttp.NewServer(
		endpoint,
		decodeOssPutObjectCallbackRequest,
		encodeOssPutObjectCallbackResponse,
		options...,
	))

	return router
}
