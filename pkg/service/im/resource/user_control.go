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
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/db"
	"kubesphere.io/im/pkg/models"
	"kubesphere.io/im/pkg/pb"
	"kubesphere.io/im/pkg/util/jsonutil"
	"kubesphere.io/im/pkg/util/strutil"
)

func CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := models.NewUser(req.Username, req.Email, req.PhoneNumber, req.Description, req.Password, req.Extra)

	// create new record
	if err := db.Global().Create(user).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, user)
		return nil, err
	}

	return &pb.CreateUserResponse{
		UserId: user.UserId,
	}, nil
}

func DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	userIds := req.UserId
	if len(userIds) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := db.Global().Begin()
	{
		tx.Delete(models.UserGroupBinding{}, `? in (?)`, constants.ColumnUserId, userIds)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		now := time.Now()
		distributes := map[string]interface{}{
			constants.ColumnStatusTime: now,
			constants.ColumnUpdateTime: now,
			constants.ColumnStatus:     constants.StatusDeleted,
		}
		if err := tx.Table(constants.TableUser).
			Where("? in (?)", constants.ColumnUserId, userIds).
			Updates(distributes).Error; err != nil {
			logger.Errorf(ctx, "Update user status failed: %+v", err)
			return nil, err
		}
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pb.DeleteUsersResponse{
		UserId: userIds,
	}, nil
}

func ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	user, err := GetUser(ctx, req.UserId)
	if err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Description != "" {
		user.Description = req.Description
	}
	if req.Email != "" {
		user.Email = strutil.SimplifyString(req.Email)
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = strutil.SimplifyString(req.PhoneNumber)
	}
	if len(req.Extra) > 0 {
		user.Extra = strutil.NewString(jsonutil.ToString(req.Extra))
	}
	user.UpdateTime = time.Now()

	if err := db.Global().Model(user).Updates(user).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pb.ModifyUserResponse{
		UserId: user.UserId,
	}, err
}

func GetUser(ctx context.Context, userId string) (*models.User, error) {
	var user = &models.User{UserId: userId}
	if err := db.Global().Model(models.User{}).Take(user).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// ignore Password
	user.Password = ""
	return user, nil
}

func GetUserWithGroup(ctx context.Context, userId string) (*models.UserWithGroup, error) {
	user, err := GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	groups, err := GetGroupsByUserIds(ctx, []string{userId})
	if err != nil {
		return nil, err
	}
	return &models.UserWithGroup{
		User:   user,
		Groups: groups,
	}, nil
}

func ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	req.GroupId = strutil.SimplifyStringList(req.GroupId)
	req.UserId = strutil.SimplifyStringList(req.UserId)
	req.Username = strutil.SimplifyStringList(req.Username)
	req.Email = strutil.SimplifyStringList(req.Email)
	req.PhoneNumber = strutil.SimplifyStringList(req.PhoneNumber)
	req.Status = strutil.SimplifyStringList(req.Status)

	limit := db.GetLimit(req.Limit)
	offset := db.GetOffset(req.Offset)

	// 1. get group users
	if len(req.GroupId) > 0 {
		userIds, err := GetUserIdsByGroupIds(ctx, req.GroupId)
		if err != nil {
			logger.Errorf(ctx, "%+v", err)
			return nil, err
		}

		if len(req.UserId) == 0 {
			req.UserId = userIds
		} else {
			var inUserIds []string
			for _, userId := range req.UserId {
				if strutil.Contains(userIds, userId) {
					inUserIds = append(inUserIds, userId)
				}
			}
			req.UserId = inUserIds
		}
	}

	var users []*models.User
	var count int

	if err := db.Global().
		AddQueryOrderDir(req, constants.ColumnCreateTime).
		BuildFilterConditions(req, constants.TableUser).
		Table(constants.TableUser).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	if err := db.Global().
		BuildFilterConditions(req, constants.TableUser).
		Table(constants.TableUser).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, user.ToPB())
	}

	return &pb.ListUsersResponse{
		UserSet: pbUsers,
		Total:   uint32(count),
	}, nil
}

func ListUsersWithGroup(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersWithGroupResponse, error) {
	response, err := ListUsers(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to list users: %+v", err)
		return nil, err
	}

	var userWithGroups []*pb.UserWithGroup
	for _, pbUser := range response.UserSet {
		groups, err := GetGroupsByUserIds(ctx, []string{pbUser.UserId})
		if err != nil {
			logger.Errorf(ctx, "Failed to get user [%s] groups: %+v", pbUser.UserId, err)
			return nil, err
		}
		var pbGroups []*pb.Group
		for _, group := range groups {
			pbGroups = append(pbGroups, group.ToPB())
		}
		userWithGroups = append(userWithGroups, &pb.UserWithGroup{
			User:     pbUser,
			GroupSet: pbGroups,
		})
	}

	return &pb.ListUsersWithGroupResponse{
		UserSet: userWithGroups,
		Total:   response.Total,
	}, nil
}
