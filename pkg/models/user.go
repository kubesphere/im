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

	"github.com/golang/protobuf/ptypes"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/pb"
	"kubesphere.io/im/pkg/util/idutil"
	"kubesphere.io/im/pkg/util/jsonutil"
	"kubesphere.io/im/pkg/util/strutil"
)

type User struct {
	UserId      string `gorm:"primary_key"`
	Username    string `gorm:"type:varchar(50);not null;unique;"`
	Email       string `gorm:"type:varchar(50);not null;unique"`
	PhoneNumber string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(1000);not null"`
	Password    string `gorm:"type:varchar(128);not null"`
	Status      string `gorm:"type:varchar(10);not null"`
	CreateTime  time.Time
	UpdateTime  time.Time
	StatusTime  time.Time
	Extra       *string `gorm:"type:JSON"`
}

type UserWithGroup struct {
	User   *User
	Groups []*Group
}

func (p *UserWithGroup) ToPB() *pb.UserWithGroup {
	var pbGroups []*pb.Group
	for _, group := range p.Groups {
		pbGroups = append(pbGroups, group.ToPB())
	}
	return &pb.UserWithGroup{
		User:     p.User.ToPB(),
		GroupSet: pbGroups,
	}
}

func NewUser(username, email, phoneNumber, description, password string, extra map[string]string) *User {
	data := jsonutil.ToString(extra)
	now := time.Now()
	user := &User{
		UserId:      idutil.GetUuid(constants.PrefixUserId),
		Username:    strutil.SimplifyString(username),
		Email:       strutil.SimplifyString(email),
		PhoneNumber: strutil.SimplifyString(phoneNumber),
		Description: description,
		Password:    password,
		Status:      constants.StatusActive,
		CreateTime:  now,
		UpdateTime:  now,
		StatusTime:  now,
		Extra:       strutil.NewString(data),
	}
	return user
}

func (p *User) ToPB() *pb.User {
	q, _ := p.ToProtoMessage()
	return q
}

func (p *User) ToProtoMessage() (*pb.User, error) {
	if p == nil {
		return new(pb.User), nil
	}
	var q = &pb.User{
		UserId:      p.UserId,
		Username:    p.Username,
		Email:       p.Email,
		PhoneNumber: p.PhoneNumber,
		Description: p.Description,
		Status:      p.Status,
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
