package object

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-sts-go-sdk/sts"
	"github.com/bluecover/qiniu_token/base"
	"github.com/bluecover/qiniu_token/model"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
)

type serviceImpl struct {
	db          *gorm.DB
	logger      kitlog.Logger
	qiniuConfig *qiniuConfig
}

// NewService creates a Object service with necessary dependencies.
func NewService(db *gorm.DB, logger kitlog.Logger, configPath string) (Service, error) {
	var qiniuViper = viper.New()
	qiniuViper.AddConfigPath(configPath)
	qiniuViper.SetConfigName("qiniu")
	if err := qiniuViper.ReadInConfig(); err != nil {
		return &serviceImpl{}, err
	}

	var qiniuConfig qiniuConfig
	if err := qiniuViper.Unmarshal(&qiniuConfig); err != nil {
		return &serviceImpl{}, err
	}

	return &serviceImpl{
		db:          db,
		logger:      logger,
		qiniuConfig: &qiniuConfig,
	}, nil
}

const (
	cloudServiceQiniu  = "qiniu"
	cloudServiceAliyun = "aliyun"
)

func decodeTokenOptions(optionsJSONString string) (map[string]interface{}, error) {
	var options map[string]interface{}
	content, err := base64.URLEncoding.DecodeString(optionsJSONString)
	if err != nil {
		return map[string]interface{}{}, err
	}
	err = json.Unmarshal(content, &options)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return options, nil
}

func getQiniuBucketFromCategory(category string) (string, error) {
	var (
		qiniuCategory2Bucket = viper.GetStringMapString("qiniu.bucket.category")
	)
	bucket, ok := qiniuCategory2Bucket[category]
	if !ok {
		return "", fmt.Errorf("unknow category %s", category)
	}
	return bucket, nil
}

func qiniuGetPrivateURL(domain string, key string, duration int) PrivateURL {
	var (
		qiniuAccessKey = viper.GetString("qiniu.access_key")
		qiniuSecretKey = viper.GetString("qiniu.secret_key")
	)
	mac := qbox.NewMac(qiniuAccessKey, qiniuSecretKey)
	deadline := time.Now().Add(time.Second * time.Duration(duration))
	url := storage.MakePrivateURL(mac, domain, key, deadline.Unix())
	return PrivateURL{
		URL:        url,
		Expiration: deadline.UTC(),
	}
}

func ossGetCredentials(bucket string, duration uint) (AccessSecrets, *base.AppError) {
	var (
		accessKeyID     = viper.GetString("aliyun.access_key_id")
		accessKeySecret = viper.GetString("aliyun.access_key_secret")
		roleArn         = viper.GetString("aliyun.role_arn_oss_wr")
		sessionName     = viper.GetString("aliyun.session_name")
	)

	stsClient := sts.NewClient(accessKeyID, accessKeySecret, roleArn, sessionName)
	resp, err := stsClient.AssumeRole(duration)
	if err != nil {
		return AccessSecrets{}, base.NewAppError(ErrAliyunSTS, errors.Wrap(err, "sts:AssumeRole"))
	}

	return AccessSecrets{
		CloudService:    cloudServiceAliyun,
		AccessKeyID:     resp.Credentials.AccessKeyId,
		AccessKeySecret: resp.Credentials.AccessKeySecret,
		Token:           resp.Credentials.SecurityToken,
		Expiration:      resp.Credentials.Expiration.UTC(),
	}, nil
}

func (impl *serviceImpl) GetUploadToken(ctx context.Context, cloud string, category string, user string) (UploadToken, *base.AppError) {
	categoryConfig, ok := impl.qiniuConfig.Category[category]
	if !ok {
		return UploadToken{}, base.NewAppError(ErrInvalidParameter, fmt.Errorf("unknown category: %s", category))
	}

	watermarkText := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("ID:%s", user)))
	// watermarkImage := base64.URLEncoding.EncodeToString([]byte("http://p6byep6mn.bkt.clouddn.com/watermark_26.png"))
	// fmt.Println(watermarkImage)

	persistentOpsList := make([]string, 0, len(categoryConfig.PersistentOps))
	for _, ops := range categoryConfig.PersistentOps {
		strFop := strings.Replace(ops.Pfop, "${wmText}", watermarkText, 1)
		saveAs := fmt.Sprintf("%s:%s", ops.SaveBucket, ops.SaveKey)
		b64SaveAs := base64.URLEncoding.EncodeToString([]byte(saveAs))
		persistentOpsList = append(persistentOpsList, fmt.Sprintf("%s|saveas/%s", strFop, b64SaveAs))
	}
	persistentOps := strings.Join(persistentOpsList, ";")

	returnKeyValues := []string{
		`"etag":$(etag)`,
		`"key":$(key)`,
		`"size":$(fsize)`,
		`"mime_type":$(mimeType)`,
	}
	for _, item := range categoryConfig.ReturnBody {
		returnKeyValues = append(returnKeyValues, item)
	}

	// returnBody := `{"etag":"$(etag)","key":"$(key)","size":$(fsize),"mimeType":$(mimeType),"persistentId":$(persistentId)}`
	returnBody := fmt.Sprintf(`{%s}`, strings.Join(returnKeyValues, ","))

	putPolicy := storage.PutPolicy{
		Scope:              categoryConfig.Scope,
		IsPrefixalScope:    int(categoryConfig.IsPrefixalScope),
		SaveKey:            categoryConfig.SaveKey,
		EndUser:            user,
		PersistentOps:      persistentOps,
		PersistentPipeline: categoryConfig.PersistentPipeline,
		Expires:            uint32(impl.qiniuConfig.TokenDuration),
		MimeLimit:          categoryConfig.MimeLimit,
		FsizeLimit:         categoryConfig.FsizeLimit,
		FsizeMin:           categoryConfig.FsizeMin,
		InsertOnly:         uint16(categoryConfig.InsertOnly),
		ReturnBody:         returnBody,
	}
	mac := qbox.NewMac(impl.qiniuConfig.AccessKey, impl.qiniuConfig.SecretKey)
	uploadToken := putPolicy.UploadToken(mac)

	return UploadToken{
		Bucket:     categoryConfig.Bucket,
		Token:      uploadToken,
		Expiration: time.Now().Add(time.Second * time.Duration(impl.qiniuConfig.TokenDuration)).UTC(),
	}, nil
}

func (impl *serviceImpl) GetAccessSecrets(ctx context.Context, cloud string, bucket string, optionsJSON string) (AccessSecrets, *base.AppError) {
	if cloud == "aliyun" {
		var (
			tokenDuration = viper.GetInt("aliyun.token_duration")
		)
		return ossGetCredentials(bucket, uint(tokenDuration))
	} else if cloud == "qiniu" {
		return AccessSecrets{
			CloudService: cloudServiceQiniu,
			Token:        "",
			Expiration:   time.Now().UTC(),
		}, nil
	} else {
		return AccessSecrets{}, base.NewAppError(ErrUnsupportedCloundService, fmt.Errorf("%s is not supported", cloud))
	}
}

func (impl *serviceImpl) GetPrivateURL(ctx context.Context, cloud string, domain string, key string) (PrivateURL, *base.AppError) {
	if cloud == "qiniu" {
		var (
			privateURLDuration = viper.GetInt("qiniu.private_url_duration")
		)
		return qiniuGetPrivateURL(domain, key, privateURLDuration), nil
	}
	return PrivateURL{}, base.NewAppError(ErrUnimplemented, fmt.Errorf("GetPrivateURL"))
}

func (impl *serviceImpl) AddObjectReference(ctx context.Context, userID uint, tag string, objInfo ObjectInfo) *base.AppError {
	obj := &model.Object{
		Cloud:       objInfo.Cloud,
		Bucket:      objInfo.Bucket,
		Key:         objInfo.Key,
		Etag:        objInfo.Etag,
		Size:        objInfo.Size,
		MimeType:    objInfo.MimeType,
		Status:      0,
		CreatedTime: time.Now(),
	}
	err := model.StoreObject(impl.db, obj)

	if err != nil && !base.IsMySQLDuplicateEntryError(err) {
		return base.NewAppError(ErrModelOperation, errors.Wrap(err, "StoreObject"))
	}

	sql := "SELECT id FROM oss WHERE cloud = ? and bucket = ? and `key` = ?"
	args := []interface{}{objInfo.Cloud, objInfo.Bucket, objInfo.Key}
	rows, err := impl.db.Raw(sql, args...).Rows()
	if err != nil {
		return base.NewAppError(ErrModelOperation, errors.Wrap(err, "AddObjectReference"))
	}
	var objectID uint
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&objectID)
		if err != nil {
			return base.NewAppError(ErrModelOperation, errors.Wrap(err, "AddObjectReference"))
		}
	}

	err = model.StoreObjectRef(impl.db, &model.ObjectRef{
		UserID:      userID,
		ObjectID:    objectID,
		Tag:         tag,
		Status:      0,
		CreatedTime: time.Now(),
	})
	if err != nil {
		return base.NewAppError(ErrModelOperation, errors.Wrap(err, "StoreObjectRef"))
	}

	return nil
}

func (impl *serviceImpl) RemoveObjectReference(ctx context.Context, userID uint, objectID uint, tag string) *base.AppError {
	err := model.DeleteObjectRef(impl.db, &model.ObjectRef{
		UserID:   userID,
		ObjectID: objectID,
		Tag:      tag,
	})
	if err != nil {
		return base.NewAppError(ErrModelOperation, errors.Wrap(err, "StoreObjectRef"))
	}

	return nil
}

func (impl *serviceImpl) GetObject(ctx context.Context, id uint) (ObjectInfo, *base.AppError) {
	mobj, err := model.FindObject(impl.db, id)
	if err != nil {
		return ObjectInfo{}, base.NewAppError(ErrNotFound, errors.Wrap(err, "FindObject"))
	}

	obj := extractModelObject(mobj)
	return *obj, nil
}

func (impl *serviceImpl) GetAllObjects(ctx context.Context) ([]ObjectInfo, *base.AppError) {
	mobjcts, err := model.FindAllObjects(impl.db)
	if err != nil {
		return make([]ObjectInfo, 0), base.NewAppError(ErrNotFound, errors.Wrap(err, "FindAllObjects"))
	}

	objInfoes := make([]ObjectInfo, 0, len(mobjcts))
	for _, mobj := range mobjcts {
		objInfoes = append(objInfoes, *extractModelObject(&mobj))
	}

	return objInfoes, nil
}

func (impl *serviceImpl) DeleteObject(ctx context.Context, id uint) *base.AppError {
	err := model.DeleteObject(impl.db, id)
	if err != nil {
		return base.NewAppError(ErrUpdateFailed, errors.Wrap(err, "DeleteObject"))
	}
	return nil
}

func extractModelObject(mobj *model.Object) *ObjectInfo {
	return &ObjectInfo{
		Cloud:    mobj.Cloud,
		Bucket:   mobj.Bucket,
		Key:      mobj.Key,
		Etag:     mobj.Etag,
		MimeType: mobj.MimeType,
		Size:     mobj.Size,
		Status:   mobj.Status,
	}
}
