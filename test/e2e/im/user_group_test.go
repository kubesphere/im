// Copyright 2019 The KubeSphere Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build integration

package im

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"kubesphere.io/im/pkg/constants"
	"kubesphere.io/im/pkg/pb"
)

func TestUserGroup(t *testing.T) {
	prepare(t)

	group := &pb.Group{
		ParentGroupId: "",
		GroupName:     "test",
		Description:   "for test",
		Extra: map[string]string{
			"num": "5",
		},
	}

	user := &pb.User{
		Username:    "test",
		Email:       "test@op.com",
		PhoneNumber: "10000000000",
		Description: "for test",
		Extra: map[string]string{
			"age": "20",
		},
	}
	password := "passw0rd"

	ctx := context.Background()

	// create parent group
	createGroupResponse, err := imClient.CreateGroup(ctx, &pb.CreateGroupRequest{
		ParentGroupId: group.ParentGroupId,
		GroupName:     group.GroupName,
		Description:   group.Description,
		Extra:         group.Extra,
	})
	require.NoError(t, err)
	group.GroupId = createGroupResponse.GroupId
	group.GroupPath = createGroupResponse.GroupId

	// create user
	createUserResponse, err := imClient.CreateUser(ctx, &pb.CreateUserRequest{
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Description: user.Description,
		Password:    password,
		Extra:       user.Extra,
	})
	require.NoError(t, err)
	user.UserId = createUserResponse.UserId

	// join group
	_, err = imClient.JoinGroup(ctx, &pb.JoinGroupRequest{
		GroupId: []string{group.GroupId},
		UserId:  []string{user.UserId},
	})
	require.NoError(t, err)

	// get group with user
	getGroupWithUserResponse, err := imClient.GetGroupWithUser(ctx, &pb.GetGroupRequest{
		GroupId: group.GroupId,
	})
	require.NoError(t, err)
	isGroupEqual(t, group, getGroupWithUserResponse.Group.Group, constants.StatusActive)
	require.EqualValues(t, len(getGroupWithUserResponse.Group.UserSet), 1)
	isUserEqual(t, user, getGroupWithUserResponse.Group.UserSet[0], constants.StatusActive)

	// get user with group
	getUserWithGroupResponse, err := imClient.GetUserWithGroup(ctx, &pb.GetUserRequest{
		UserId: user.UserId,
	})
	require.NoError(t, err)
	isUserEqual(t, user, getUserWithGroupResponse.User.User, constants.StatusActive)
	require.EqualValues(t, len(getUserWithGroupResponse.User.GroupSet), 1)
	isGroupEqual(t, group, getUserWithGroupResponse.User.GroupSet[0], constants.StatusActive)

	// list groups with user
	listGroupsWithUserResponse, err := imClient.ListGroupsWithUser(ctx, &pb.ListGroupsRequest{
		GroupId: []string{group.GroupId},
		Status:  []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listGroupsWithUserResponse.Total, 1)
	require.EqualValues(t, len(listGroupsWithUserResponse.GroupSet[0].UserSet), 1)
	isGroupEqual(t, group, listGroupsWithUserResponse.GroupSet[0].Group, constants.StatusActive)
	isUserEqual(t, user, listGroupsWithUserResponse.GroupSet[0].UserSet[0], constants.StatusActive)

	// list users with group
	listUsersWithGroupResponse, err := imClient.ListUsersWithGroup(ctx, &pb.ListUsersRequest{
		GroupId: []string{group.GroupId},
		Status:  []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listUsersWithGroupResponse.Total, 1)
	require.EqualValues(t, len(listUsersWithGroupResponse.UserSet[0].GroupSet), 1)
	isGroupEqual(t, group, listUsersWithGroupResponse.UserSet[0].GroupSet[0], constants.StatusActive)
	isUserEqual(t, user, listUsersWithGroupResponse.UserSet[0].User, constants.StatusActive)

	// list users with root group id
	listUsersWithGroupResponse, err = imClient.ListUsersWithGroup(ctx, &pb.ListUsersRequest{
		RootGroupId: []string{group.GroupId},
		Status:      []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listUsersWithGroupResponse.Total, 1)
	require.EqualValues(t, len(listUsersWithGroupResponse.UserSet[0].GroupSet), 1)
	isGroupEqual(t, group, listUsersWithGroupResponse.UserSet[0].GroupSet[0], constants.StatusActive)
	isUserEqual(t, user, listUsersWithGroupResponse.UserSet[0].User, constants.StatusActive)

	// delete group, has user, can not delete
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{group.GroupId},
	})
	require.Error(t, err)

	// leave group
	_, err = imClient.LeaveGroup(ctx, &pb.LeaveGroupRequest{
		GroupId: []string{group.GroupId},
		UserId:  []string{user.UserId},
	})
	require.NoError(t, err)

	// delete group
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{group.GroupId},
	})
	require.NoError(t, err)

	// delete user
	_, err = imClient.DeleteUsers(ctx, &pb.DeleteUsersRequest{
		UserId: []string{user.UserId},
	})
	require.NoError(t, err)

}
