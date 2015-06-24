// Copyright 2015 Lorenzo Landolfi.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package resourcer

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type GormHandlers interface {
	Resource(GormHandlers, *gorm.DB, gin.Params) (interface{}, error)
	ResourceClass() (reflect.Type, string)
}

// Gets the current resource from the ID in the URL
type ResourceGorm struct{}

// Gets ALL the resources present in db for this model
type CollectionGorm struct{}

// Implementation of Resource for ResourceGorm
// If you want to create your own just don't do embedding
// of ResourceGorm struct into your GormHandlers and
// implement your own (Or create your own resource getter
func (_ *ResourceGorm) Resource(h GormHandlers, db *gorm.DB, p gin.Params) (interface{}, error) {
	rtype, rname := h.ResourceClass()
	id, ok := p.Get(rname + "_id")
	if ok != true {
		return nil, nil
	}

	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	res := reflect.New(rtype).Interface()
	if err := db.Find(res, rid).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// Same of ResourceGorm but for collection, see above
func (_ *CollectionGorm) Resource(h GormHandlers, db *gorm.DB, p gin.Params) (interface{}, error) {
	rtype, _ := h.ResourceClass()

	res := reflect.New(reflect.SliceOf(rtype)).Interface()
	if err := db.Find(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func ResourcerGorm(h GormHandlers, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gets resource and sets it into Context
		res, err := h.Resource(h, db, c.Params)

		if err != nil {
			if gin.IsDebugging() {
				c.AbortWithError(404, err)
			} else {
				c.AbortWithStatus(404)
			}
		}
		c.Resource = res
		// Before request
		c.Next()
		// After request
	}
}
