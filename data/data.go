// Copyright 2015 mint.zhao.chiu@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package data

import (
	"github.com/gin-gonic/gin"
)

var (
	DefaultKey = "github.com/gin-gonic/contrib/data"
)

type DataRender interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	Del(key string)
	Render() gin.H

    Error(err string)
}

type data struct {
	data gin.H
}

func (this *data) Get(key string) interface{} {
	return this.data[key]
}

func (this *data) Set(key string, val interface{}) {
	this.data[key] = val
}

func (this *data) Del(key string) {
	delete(this.data, key)
}

func (this *data) Render() gin.H {
	return this.data
}

func (this *data) Error(err string) {
    this.Set("error", err)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := &data{data: make(gin.H)}
		c.Set(DefaultKey, data)
		c.Next()
	}
}

func Data(ctx *gin.Context) DataRender {
	return ctx.MustGet(DefaultKey).(DataRender)
}
