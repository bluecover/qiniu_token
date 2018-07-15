package callback

// Callback service implementation

import (
	"context"
	"fmt"
	"time"

	"github.com/bluecover/qiniu_token/base"
	"github.com/bluecover/qiniu_token/model"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewService(db *gorm.DB) Service {
	return &serviceImpl{
		db: db,
	}
}

type serviceImpl struct {
	db     *gorm.DB
	logger log.Logger
}

func (impl *serviceImpl) OssPutObjectCallback(ctx context.Context, param OssCallbackParam) (CallbackResult, *base.AppError) {
	verifyUserToken(param.AppUserID, param.AppUserToken)

	objFilename := createFilename(param.AppBusiness, param.Etag, param.ImageFormat)

	if len(param.OriginName) > 32 {
		param.OriginName = param.OriginName[0:32]
	}

	obj := &model.Object{
		Cloud:       "aliyun",
		Bucket:      param.Bucket,
		Key:         objFilename,
		Etag:        param.Etag,
		Size:        uint(param.Size),
		Status:      0,
		CreatedTime: time.Now(),
	}
	err := model.StoreObject(impl.db, obj)
	if err != nil {
		if base.IsMySQLDuplicateEntryError(err) {
			return CallbackResult{}, base.NewAppError(ErrAlreadyExists, errors.Wrap(err, "StoreObject"))
		} else {
			return CallbackResult{}, base.NewAppError(ErrModelFunctionFailed, errors.Wrap(err, "StoreObject"))
		}
	}

	return CallbackResult{ObjID: obj.ID, ObjFilename: obj.Key}, nil
}

func verifyUserToken(userID uint, token string) bool {
	// TODO: Verify user token (param.token)
	return true
}

func createFilename(business string, etag string, format string) string {
	yearMonthStr := time.Now().Format("2006/01")
	return fmt.Sprintf("%s/%s/%s.%s", business, yearMonthStr, etag, format)
}
