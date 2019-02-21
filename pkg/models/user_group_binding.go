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

package models

import (
	"time"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/util/idutil"
)

type UserGroupBinding struct {
	Id         string    `gorm:"type:varchar(50);primary_key"`
	GroupId    string    `gorm:"type:varchar(50);not null"`
	UserId     string    `gorm:"type:varchar(50);not null"`
	CreateTime time.Time `gorm:"default CURRENT_TIMESTAMP"`
}

func NewUserGroupBinding(userId, groupId string) *UserGroupBinding {
	return &UserGroupBinding{
		Id:         idutil.GetUuid(constants.PrefixUserGroupBindingId),
		GroupId:    groupId,
		UserId:     userId,
		CreateTime: time.Now(),
	}
}
