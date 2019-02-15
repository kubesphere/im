/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"strings"
	"sync"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/gorm"
	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/util/strutil"
)

const (
	DefaultSelectLimit = 200
)

func GetLimit(n uint32) uint32 {
	if n < 0 {
		n = 0
	}
	if n > DefaultSelectLimit {
		n = DefaultSelectLimit
	}
	return n
}

func GetOffset(n uint32) uint32 {
	if n < 0 {
		n = 0
	}
	return n
}

var global *Database
var globalMutex sync.RWMutex

func SetGlobal(database *Database) {
	globalMutex.Lock()
	global = database
	globalMutex.Unlock()
}

func Global() *Database {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return global
}

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}
type RequestWithSortKey interface {
	Request
	GetSortKey() string
}
type RequestWithReverse interface {
	RequestWithSortKey
	GetReverse() bool
}

const (
	TagName              = "json"
	SearchWordColumnName = "search_word"
)

func getReqValue(param interface{}) interface{} {
	switch value := param.(type) {
	case string:
		if value == "" {
			return nil
		}
		return []string{value}
	case *wrappers.StringValue:
		if value == nil {
			return nil
		}
		return []string{value.GetValue()}
	case *wrappers.Int32Value:
		if value == nil {
			return nil
		}
		return []int32{value.GetValue()}
	case []string:
		var values []string
		for _, v := range value {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) == 0 {
			return nil
		}
		return values
	}
	return nil
}

func GetDisplayColumns(displayColumns []string, wholeColumns []string) []string {
	if displayColumns == nil {
		return wholeColumns
	} else if len(displayColumns) == 0 {
		return nil
	} else {
		var newDisplayColumns []string
		for _, column := range displayColumns {
			if strutil.Contains(wholeColumns, column) {
				newDisplayColumns = append(newDisplayColumns, column)
			}
		}
		return newDisplayColumns
	}
}

func getFieldName(field *structs.Field) string {
	tag := field.Tag(TagName)
	t := strings.Split(tag, ",")
	if len(t) == 0 {
		return "-"
	}
	return t[0]
}

type Chain struct {
	*gorm.DB
}

func GetChain(tx *gorm.DB) *Chain {
	return &Chain{
		tx,
	}
}

func (c *Chain) BuildFilterConditions(req Request, tableName string, exclude ...string) *Chain {
	return c.buildFilterConditions(req, tableName, exclude...)
}

func (c *Chain) getSearchFilter(tableName string, value interface{}, exclude ...string) {
	var conditions []string
	if vs, ok := value.([]string); ok {
		for _, v := range vs {
			for _, column := range constants.SearchColumns[tableName] {
				if strutil.Contains(exclude, column) {
					continue
				}
				// if column suffix is _id, must exact match
				if strings.HasSuffix(column, "_id") {
					conditions = append(conditions, column+" = '"+v+"'")
				} else {
					likeV := "%" + strutil.SimplifyString(v) + "%"
					conditions = append(conditions, column+" LIKE '"+likeV+"'")
				}
			}
		}
	} else if value != nil {
		logger.Warnf(nil, "search_word [%+v] is not string", value)
	}
	condition := strings.Join(conditions, " OR ")
	c.DB = c.DB.Where(condition)
}

func (c *Chain) buildFilterConditions(req Request, tableName string, exclude ...string) *Chain {
	for _, field := range structs.Fields(req) {
		column := getFieldName(field)
		param := field.Value()
		indexedColumns, ok := constants.IndexedColumns[tableName]
		if ok && strutil.Contains(indexedColumns, column) {
			value := getReqValue(param)
			if value != nil {
				key := column
				c.DB = c.Where(key+" in (?)", value)
			}
		}
		if column == SearchWordColumnName && strutil.Contains(constants.SearchWordColumnTable, tableName) {
			value := getReqValue(param)
			c.getSearchFilter(tableName, value, exclude...)
		}
	}
	return c
}

func (c *Chain) AddQueryOrderDir(req Request, defaultColumn string) *Chain {
	order := "DESC"
	if r, ok := req.(RequestWithReverse); ok {
		if r.GetReverse() {
			order = "ASC"
		}
	}
	if r, ok := req.(RequestWithSortKey); ok {
		s := r.GetSortKey()
		if s != "" {
			defaultColumn = s
		}
	}
	c.DB = c.Order(defaultColumn + " " + order)
	return c
}
