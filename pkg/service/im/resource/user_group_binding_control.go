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

package resource

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/db"
	"kubesphere.io/im/pkg/models"
	"kubesphere.io/im/pkg/pb"
)

func GetUserGroupBindings(ctx context.Context, groupId, userIds []string) ([]*models.UserGroupBinding, error) {
	var userGroupBindings []*models.UserGroupBinding
	if err := db.Global().Model(models.User{}).
		Where("? in (?) and ? in (?)", constants.ColumnGroupId, groupId, constants.ColumnUserId, userIds).
		Find(&userGroupBindings).
		Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return userGroupBindings, nil
}

func JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	if len(req.UserId) == 0 || len(req.GroupId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user id or group id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// check user in group
	userGroupBindings, err := GetUserGroupBindings(ctx, req.GroupId, req.UserId)
	if err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}
	if len(userGroupBindings) != 0 {
		err := status.Errorf(codes.PermissionDenied, "user already in group")
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	tx := db.Global().Begin()
	{
		for _, groupId := range req.GroupId {
			for _, userId := range req.UserId {
				if err := tx.Create(models.NewUserGroupBinding(groupId, userId)).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pb.JoinGroupResponse{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	}, nil
}

func LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	if len(req.UserId) == 0 || len(req.GroupId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user id or group id")
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	// check user in group
	userGroupBindings, err := GetUserGroupBindings(ctx, req.GroupId, req.UserId)
	if err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}
	if len(userGroupBindings) != len(req.UserId)*len(req.GroupId) {
		err := status.Errorf(codes.PermissionDenied, "user not in group")
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	if err := db.Global().Delete(models.UserGroupBinding{},
		`user_id in (?) and group_id in (?)`,
		req.UserId, req.GroupId,
	).Error; err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	return &pb.LeaveGroupResponse{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	}, nil
}

func GetGroupsByUserIds(ctx context.Context, userIds []string) ([]*models.Group, error) {
	const query = `
		select user_group.* from
			user, user_group, user_group_binding
		where
			user_group_binding.user_id=user.user_id and
			user_group_binding.group_id=user_group.group_id and
			user.user_id in (?)
	`
	var groups []*models.Group
	if err := db.Global().Raw(query, userIds).Scan(&groups).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return groups, nil
}

func GetUsersByGroupIds(ctx context.Context, groupIds []string) ([]*models.User, error) {
	const query = `
		select user.* from
			user, user_group, user_group_binding
		where
			user_group_binding.user_id=user.user_id and
			user_group_binding.group_id=user_group.group_id and
 			user_group.group_id in (?)
	`
	var users []*models.User
	if err := db.Global().Raw(query, groupIds).Scan(users).Error; err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	return users, nil
}

func GetUserIdsByGroupIds(ctx context.Context, groupIds []string) ([]string, error) {
	rows, err := db.Global().Table(constants.TableUserGroupBinding).
		Select(constants.ColumnUserId).
		Where("? in (?)", constants.ColumnGroupId, groupIds).
		Rows()
	if err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}
	var userIds []string
	for rows.Next() {
		var userId string
		rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}
