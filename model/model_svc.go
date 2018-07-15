package model

import (
	"github.com/bluecover/qiniu_token/base"
	"github.com/jinzhu/gorm"
)

// DB record status
const (
	StatusNormal  = 0
	StatusDeleted = -1
)

// FindObject retrieves the Object specified by id.
func FindObject(db *gorm.DB, id uint) (*Object, error) {
	obj := new(Object)
	err := db.Where(&Object{ID: id}).First(obj).Error
	if err != nil {
		return &Object{}, err
	}
	return obj, nil
}

// FindAllObjects retrieves all objects.
func FindAllObjects(db *gorm.DB) ([]Object, error) {
	var objs []Object
	err := db.Find(&objs).Error
	return objs, err
}

// StoreObject create the new Object record.
func StoreObject(db *gorm.DB, obj *Object) error {
	result := db.Create(obj)
	return result.Error
}

// DeleteObject executes soft-delete on the Object specified by id.
func DeleteObject(db *gorm.DB, id uint) error {
	err := db.Model(&Object{ID: id}).UpdateColumn("status", StatusDeleted).Error
	return err
}

// StoreObjectRef executes creating operations for a new ObjectRef.
func StoreObjectRef(db *gorm.DB, objRef *ObjectRef) error {
	err := db.Create(objRef).Error
	if err != nil && !base.IsMySQLDuplicateEntryError(err) {
		return err
	}
	objRef.Status = StatusDeleted
	err = db.Model(objRef).UpdateColumn("status", StatusNormal).Error
	return err
}

// DeleteObjectRef executes soft-delete on the ObjectRef specified by id.
func DeleteObjectRef(db *gorm.DB, objRef *ObjectRef) error {
	err := db.Model(objRef).UpdateColumn("status", StatusDeleted).Error
	return err
}

// AllModels retrieve a list of all model objects with empty values.
func AllModels() []interface{} {
	return []interface{}{
		Object{},
		ObjectRef{},
	}
}
