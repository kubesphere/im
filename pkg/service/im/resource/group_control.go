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
	"strings"
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

func CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	parentGroupId := strutil.SimplifyString(req.ParentGroupId)
	parentGroupPath, err := GetParentGroupPath(ctx, parentGroupId)
	if err != nil {
		return nil, err
	}

	group := models.NewGroup(parentGroupId, parentGroupPath, req.GroupName, req.Description, req.Extra)

	var allParentGroupIds []string
	// skip groupId
	for _, groupId := range strings.Split(group.GroupPath, ".") {
		if groupId != group.GroupId {
			allParentGroupIds = append(allParentGroupIds, groupId)
		}
	}

	// check all parent group id exists
	if len(allParentGroupIds) > 0 {

		var total int
		db.Global().Table(constants.TableGroup).
			Where(constants.ColumnGroupId+" in (?)", allParentGroupIds).
			Count(&total)

		if err := db.Global().Error; err != nil {
			logger.Errorf(ctx, "Get parent group ids failed: %+v", err)
			return nil, err
		}
		if total != len(allParentGroupIds) {
			err := status.Errorf(codes.InvalidArgument,
				"some groupId in allParentGroupIds (%q) do not exists",
				group.GroupPath,
			)
			logger.Errorf(ctx, "%+v", err)
			return nil, err
		}
	}

	// create new record
	if err := db.Global().Create(group).Error; err != nil {
		logger.Errorf(ctx, "Insert group failed: %+v", err)
		return nil, err
	}

	return &pb.CreateGroupResponse{
		GroupId: group.GroupId,
	}, nil
}

func DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	groupIds := req.GroupId
	if len(groupIds) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty group id")
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	// 1. check sub groups
	subGroupIds, err := getAllSubGroupIds(ctx, groupIds)
	if err != nil {
		return nil, err
	}
	if len(subGroupIds) > 0 {
		err := status.Errorf(codes.PermissionDenied, "there are still sub groups in group: %v", subGroupIds)
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	// 2. check users
	users, err := GetUserIdsByGroupIds(ctx, groupIds)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		err := status.Errorf(codes.PermissionDenied, "there are still users in group: %v", subGroupIds)
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	// 3. update group status to deleted
	now := time.Now()
	attributes := map[string]interface{}{
		constants.ColumnStatusTime: now,
		constants.ColumnUpdateTime: now,
		constants.ColumnStatus:     constants.StatusDeleted,
	}
	if err := db.Global().Table(constants.TableGroup).
		Where(constants.ColumnGroupId+" in (?)", groupIds).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update group status failed: %+v", err)
		return nil, err
	}

	return &pb.DeleteGroupsResponse{
		GroupId: groupIds,
	}, nil
}

func ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	groupId := req.GroupId
	group, err := GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})
	if req.ParentGroupId != "" && req.ParentGroupId != group.ParentGroupId {
		group.ParentGroupId = req.ParentGroupId
		parentGroupPath, err := GetParentGroupPath(ctx, group.ParentGroupId)
		if err != nil {
			return nil, err
		}
		groupPath := models.GetGroupPath(parentGroupPath, group.GroupId)
		attributes[constants.ColumnParentGroupId] = req.ParentGroupId
		attributes[constants.ColumnGroupPath] = groupPath
		attributes[constants.ColumnGroupPathLevel] = strings.Count(strutil.SimplifyString(groupPath), constants.GroupPathSep) + 1
	}
	if req.GroupName != "" {
		attributes[constants.ColumnGroupName] = req.GroupName
	}
	if req.Description != "" {
		attributes[constants.ColumnDescription] = req.Description
	}
	if len(req.Extra) > 0 {
		attributes[constants.ColumnExtra] = strutil.NewString(jsonutil.ToString(req.Extra))
	}
	attributes[constants.ColumnUpdateTime] = time.Now()

	if err := db.Global().Table(constants.TableGroup).
		Where(constants.ColumnGroupId+" = ?", groupId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update group [%s] failed: %+v", groupId, err)
		return nil, err
	}

	return &pb.ModifyGroupResponse{
		GroupId: groupId,
	}, nil
}

func GetParentGroupPath(ctx context.Context, parentGroupId string) (string, error) {
	parentGroupPath := ""
	if parentGroupId != "" {
		parentGroup, err := GetGroup(ctx, parentGroupId)
		if err != nil {
			err = status.Errorf(codes.InvalidArgument, "get parent group failed: %v", err)
			logger.Errorf(ctx, "%+v", err)
			return parentGroupPath, err
		}
		parentGroupPath = parentGroup.GroupPath
	}
	return parentGroupPath, nil
}

func GetGroup(ctx context.Context, groupId string) (*models.Group, error) {
	var group = &models.Group{GroupId: groupId}
	if err := db.Global().Table(constants.TableGroup).Take(group).Error; err != nil {
		logger.Errorf(ctx, "Get group [%s] failed: %+v", groupId, err)
		return nil, err
	}

	return group, nil
}

func GetGroupWithUser(ctx context.Context, groupId string) (*models.GroupWithUser, error) {
	group, err := GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}
	users, err := GetUsersByGroupIds(ctx, []string{groupId})
	if err != nil {
		return nil, err
	}
	return &models.GroupWithUser{
		Group: group,
		Users: users,
	}, nil
}

func ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	req.ParentGroupId = strutil.SimplifyStringList(req.ParentGroupId)
	req.GroupId = strutil.SimplifyStringList(req.GroupId)
	req.GroupPath = strutil.SimplifyStringList(req.GroupPath)
	req.GroupName = strutil.SimplifyStringList(req.GroupName)
	req.Status = strutil.SimplifyStringList(req.Status)

	limit := db.GetLimitFromRequest(req)
	offset := db.GetOffsetFromRequest(req)

	var groups []*models.Group
	var count int

	if err := db.GetChain(db.Global().Table(constants.TableGroup)).
		AddQueryOrderDir(req, constants.ColumnCreateTime).
		BuildFilterConditions(req, constants.TableGroup).
		Offset(offset).
		Limit(limit).
		Find(&groups).Error; err != nil {
		logger.Errorf(ctx, "List group failed: %+v", err)
		return nil, err
	}

	if err := db.GetChain(db.Global().Table(constants.TableGroup)).
		BuildFilterConditions(req, constants.TableGroup).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "List group count failed: %+v", err)
		return nil, err
	}

	var pbGroups []*pb.Group
	for _, group := range groups {
		pbGroups = append(pbGroups, group.ToPB())
	}

	return &pb.ListGroupsResponse{
		GroupSet: pbGroups,
		Total:    uint32(count),
	}, nil
}

func ListGroupsWithUser(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsWithUserResponse, error) {
	response, err := ListGroups(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "List groups failed: %+v", err)
		return nil, err
	}

	var groupWithUsers []*pb.GroupWithUser
	for _, pbGroup := range response.GroupSet {
		users, err := GetUsersByGroupIds(ctx, []string{pbGroup.GroupId})
		if err != nil {
			return nil, err
		}
		var pbUsers []*pb.User
		for _, user := range users {
			pbUsers = append(pbUsers, user.ToPB())
		}
		groupWithUsers = append(groupWithUsers, &pb.GroupWithUser{
			Group:   pbGroup,
			UserSet: pbUsers,
		})
	}

	return &pb.ListGroupsWithUserResponse{
		GroupSet: groupWithUsers,
		Total:    response.Total,
	}, nil
}

func getAllSubGroupIds(ctx context.Context, groupIds []string) ([]string, error) {
	var groups []*models.Group

	tx := db.Global().Table(constants.TableGroup)
	for _, groupId := range groupIds {
		likeGroupId := "%" + groupId + "%"
		tx = tx.Or(constants.ColumnGroupPath+" LIKE ?", likeGroupId)
	}
	if err := tx.Find(&groups).Error; err != nil {
		logger.Errorf(ctx, "Get all sub groups failed: %+v", err)
		return nil, err
	}

	var allGroupId []string
	for _, group := range groups {
		if !strutil.Contains(groupIds, group.GroupId) {
			allGroupId = append(allGroupId, group.GroupId)
		}
	}

	return allGroupId, nil
}
