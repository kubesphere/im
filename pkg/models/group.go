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
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/pb"
	"kubesphere.io/im/pkg/util/idutil"
	"kubesphere.io/im/pkg/util/jsonutil"
	"kubesphere.io/im/pkg/util/strutil"
)

type Group struct {
	ParentGroupId string `gorm:"type:varchar(50);not null"`
	GroupId       string `gorm:"primary_key"`
	GroupPath     string `gorm:"type:varchar(255);not null;unique"`
	GroupName     string `gorm:"type:varchar(50);not null"`
	Description   string `gorm:"type:varchar(1000);not null"`
	Status        string `gorm:"type:varchar(50);not null"`
	CreateTime    time.Time
	UpdateTime    time.Time
	StatusTime    time.Time
	Extra         *string `gorm:"type:JSON"`

	// internal
	GroupPathLevel int
}

type GroupWithUser struct {
	Group *Group
	Users []*User
}

func (p *GroupWithUser) ToPB() *pb.GroupWithUser {
	var pbUsers []*pb.User
	for _, user := range p.Users {
		pbUsers = append(pbUsers, user.ToPB())
	}
	return &pb.GroupWithUser{
		Group:   p.Group.ToPB(),
		UserSet: pbUsers,
	}
}

func GetGroupPath(parentGroupPath, groupId string) string {
	var groupPath string
	if parentGroupPath != "" {
		groupPath = parentGroupPath + "." + groupId
	} else {
		groupPath = groupId
	}
	return groupPath
}

func NewGroup(parentGroupId, parentGroupPath, groupName, description string, extra map[string]string) *Group {
	groupId := idutil.GetUuid(constants.PrefixGroupId)
	groupPath := GetGroupPath(parentGroupPath, groupId)
	data := jsonutil.ToString(extra)
	now := time.Now()
	group := &Group{
		ParentGroupId:  strutil.SimplifyString(parentGroupId),
		GroupId:        strutil.SimplifyString(groupId),
		GroupPath:      strutil.SimplifyString(groupPath),
		GroupName:      groupName,
		Description:    description,
		Status:         constants.StatusActive,
		CreateTime:     now,
		UpdateTime:     now,
		StatusTime:     now,
		Extra:          strutil.NewString(data),
		GroupPathLevel: strings.Count(strutil.SimplifyString(groupPath), constants.GroupPathSep) + 1,
	}
	return group
}

func (p *Group) ToProtoMessage() (*pb.Group, error) {
	if p == nil {
		return new(pb.Group), nil
	}
	var q = &pb.Group{
		ParentGroupId: p.ParentGroupId,
		GroupId:       p.GroupId,
		GroupPath:     p.GroupPath,
		GroupName:     p.GroupName,
		Description:   p.Description,
		Status:        p.Status,
	}

	q.CreateTime, _ = ptypes.TimestampProto(p.CreateTime)
	q.UpdateTime, _ = ptypes.TimestampProto(p.UpdateTime)
	q.StatusTime, _ = ptypes.TimestampProto(p.StatusTime)

	if p.Extra != nil && *p.Extra != "" {
		if q.Extra == nil {
			q.Extra = make(map[string]string)
		}
		err := jsonutil.Decode([]byte(*p.Extra), &q.Extra)
		if err != nil {
			return q, err
		}
	}
	return q, nil
}

func (p *Group) ToPB() *pb.Group {
	q, _ := p.ToProtoMessage()
	return q
}

func (p *Group) IsValidSortKey(key string) bool {
	var validKeys = []string{
		"parent_group_id",
		"group_id",
		"group_path",
		"group_name",
		"description",
		"status",
		"create_time",
		"update_time",
		"status_time",
	}
	for _, k := range validKeys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}
