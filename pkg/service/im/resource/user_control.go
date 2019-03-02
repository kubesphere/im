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
	"kubesphere.io/im/pkg/global"
	"kubesphere.io/im/pkg/models"
	"kubesphere.io/im/pkg/pb"
	"kubesphere.io/im/pkg/util/jsonutil"
	"kubesphere.io/im/pkg/util/stringutil"
)

func CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := models.NewUser(req.Username, req.Email, req.PhoneNumber, req.Description, req.Password, req.Extra)

	// create new record
	if err := global.Global().Database.Create(user).Error; err != nil {
		logger.Errorf(ctx, "Insert user failed: %+v", err)
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
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	tx := global.Global().Database.Begin()
	{
		tx.Delete(models.UserGroupBinding{}, constants.ColumnUserId+" in (?)", userIds)
		if err := tx.Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Delete user group binding failed: %+v", err)
			return nil, err
		}

		now := time.Now()
		attributes := map[string]interface{}{
			constants.ColumnStatusTime: now,
			constants.ColumnUpdateTime: now,
			constants.ColumnStatus:     constants.StatusDeleted,
		}
		if err := tx.Table(constants.TableUser).
			Where(constants.ColumnUserId+" in (?)", userIds).
			Updates(attributes).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Update user status failed: %+v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "Delete user failed: %+v", err)
		return nil, err
	}

	return &pb.DeleteUsersResponse{
		UserId: userIds,
	}, nil
}

func ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	userId := req.UserId
	_, err := GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})
	if req.Username != "" {
		attributes[constants.ColumnUsername] = req.Username
	}
	if req.Description != "" {
		attributes[constants.ColumnDescription] = req.Description
	}
	if req.Email != "" {
		attributes[constants.ColumnEmail] = stringutil.SimplifyString(req.Email)
	}
	if req.PhoneNumber != "" {
		attributes[constants.ColumnPhoneNumber] = stringutil.SimplifyString(req.PhoneNumber)
	}
	if len(req.Extra) > 0 {
		attributes[constants.ColumnExtra] = stringutil.NewString(jsonutil.ToString(req.Extra))
	}
	attributes[constants.ColumnUpdateTime] = time.Now()

	if err := global.Global().Database.Table(constants.TableUser).
		Where(constants.ColumnUserId+" = ?", userId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update user [%s] failed: %+v", userId, err)
		return nil, err
	}

	return &pb.ModifyUserResponse{
		UserId: userId,
	}, err
}

func GetUser(ctx context.Context, userId string) (*models.User, error) {
	var user = &models.User{UserId: userId}
	if err := global.Global().Database.Table(constants.TableUser).
		Take(user).Error; err != nil {
		logger.Errorf(ctx, "Get user [%s] failed: %+v", userId, err)
		return nil, err
	}

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
	req.GroupId = stringutil.SimplifyStringList(req.GroupId)
	req.UserId = stringutil.SimplifyStringList(req.UserId)
	req.Username = stringutil.SimplifyStringList(req.Username)
	req.Email = stringutil.SimplifyStringList(req.Email)
	req.PhoneNumber = stringutil.SimplifyStringList(req.PhoneNumber)
	req.Status = stringutil.SimplifyStringList(req.Status)

	limit := db.GetLimitFromRequest(req)
	offset := db.GetOffsetFromRequest(req)

	var pbUsers []*pb.User

	// get group
	if len(req.RootGroupId) > 0 {
		allGroupIds, err := getAllSubGroupIds(ctx, req.RootGroupId)
		if err != nil {
			return nil, err
		}
		allGroupIds = append(allGroupIds, req.RootGroupId...)

		if len(req.GroupId) == 0 {
			req.GroupId = allGroupIds
		} else {
			var inGroupIds []string
			for _, groupId := range req.GroupId {
				if stringutil.Contains(allGroupIds, groupId) {
					inGroupIds = append(inGroupIds, groupId)
				}
			}
			req.GroupId = inGroupIds
		}
		if len(req.GroupId) == 0 {
			return &pb.ListUsersResponse{
				UserSet: pbUsers,
				Total:   0,
			}, nil
		}
	}

	// get group users
	if len(req.GroupId) > 0 {
		userIds, err := GetUserIdsByGroupIds(ctx, req.GroupId)
		if err != nil {
			return nil, err
		}

		if len(req.UserId) == 0 {
			req.UserId = userIds
		} else {
			var inUserIds []string
			for _, userId := range req.UserId {
				if stringutil.Contains(userIds, userId) {
					inUserIds = append(inUserIds, userId)
				}
			}
			req.UserId = inUserIds
		}
		if len(req.UserId) == 0 {
			return &pb.ListUsersResponse{
				UserSet: pbUsers,
				Total:   0,
			}, nil
		}
	}

	var users []*models.User
	var count int

	if err := db.GetChain(global.Global().Database.Table(constants.TableUser)).
		AddQueryOrderDir(req, constants.ColumnCreateTime).
		BuildFilterConditions(req, constants.TableUser).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		logger.Errorf(ctx, "List users failed: %+v", err)
		return nil, err
	}

	if err := db.GetChain(global.Global().Database.Table(constants.TableUser)).
		BuildFilterConditions(req, constants.TableUser).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "List users count failed: %+v", err)
		return nil, err
	}

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
		logger.Errorf(ctx, "List users failed: %+v", err)
		return nil, err
	}

	var userWithGroups []*pb.UserWithGroup
	for _, pbUser := range response.UserSet {
		groups, err := GetGroupsByUserIds(ctx, []string{pbUser.UserId})
		if err != nil {
			logger.Errorf(ctx, "Get user [%s] groups failed: %+v", pbUser.UserId, err)
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
