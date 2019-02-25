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

func isGroupEqual(t *testing.T, oldGroup, newGroup *pb.Group, status string) bool {
	require.Equal(t, oldGroup.ParentGroupId, newGroup.ParentGroupId)
	require.Equal(t, oldGroup.GroupPath, newGroup.GroupPath)
	require.Equal(t, oldGroup.GroupName, newGroup.GroupName)
	require.Equal(t, oldGroup.Description, newGroup.Description)
	require.Equal(t, oldGroup.Extra["num"], newGroup.Extra["num"])
	require.Equal(t, status, newGroup.Status)
	return true
}

func TestGroup(t *testing.T) {
	prepare(t)

	parentGroup := &pb.Group{
		ParentGroupId: "",
		GroupName:     "test parent",
		Description:   "for test parent",
		Extra: map[string]string{
			"num": "5",
		},
	}

	childGroup1 := &pb.Group{
		GroupName:   "test child1",
		Description: "for test child1",
		Extra: map[string]string{
			"num": "10",
		},
	}

	childGroup2 := &pb.Group{
		GroupName:   "test child2",
		Description: "for test child2",
		Extra: map[string]string{
			"num": "15",
		},
	}

	ctx := context.Background()

	// create parent group
	createGroupResponse, err := imClient.CreateGroup(ctx, &pb.CreateGroupRequest{
		ParentGroupId: parentGroup.ParentGroupId,
		GroupName:     parentGroup.GroupName,
		Description:   parentGroup.Description,
		Extra:         parentGroup.Extra,
	})
	require.NoError(t, err)
	parentGroup.GroupId = createGroupResponse.GroupId
	parentGroup.GroupPath = createGroupResponse.GroupId

	// get parent group
	getGroupResponse, err := imClient.GetGroup(ctx, &pb.GetGroupRequest{
		GroupId: parentGroup.GroupId,
	})
	require.NoError(t, err)
	isGroupEqual(t, parentGroup, getGroupResponse.Group, constants.StatusActive)

	// list parent group, use group id
	listGroupsResponse, err := imClient.ListGroups(ctx, &pb.ListGroupsRequest{
		GroupId: []string{parentGroup.GroupId},
		Status:  []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listGroupsResponse.Total, 1)
	isGroupEqual(t, parentGroup, listGroupsResponse.GroupSet[0], constants.StatusActive)

	// create child group1
	createGroupResponse, err = imClient.CreateGroup(ctx, &pb.CreateGroupRequest{
		ParentGroupId: parentGroup.GroupId,
		GroupName:     childGroup1.GroupName,
		Description:   childGroup1.Description,
		Extra:         childGroup1.Extra,
	})
	require.NoError(t, err)
	childGroup1.GroupId = createGroupResponse.GroupId
	childGroup1.ParentGroupId = parentGroup.GroupId
	childGroup1.GroupPath = childGroup1.ParentGroupId + "." + childGroup1.GroupId

	// list child group1, use group path
	listGroupsResponse, err = imClient.ListGroups(ctx, &pb.ListGroupsRequest{
		GroupPath: []string{childGroup1.GroupPath},
		Status:    []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listGroupsResponse.Total, 1)
	isGroupEqual(t, childGroup1, listGroupsResponse.GroupSet[0], constants.StatusActive)

	// create child group2
	createGroupResponse, err = imClient.CreateGroup(ctx, &pb.CreateGroupRequest{
		ParentGroupId: parentGroup.GroupId,
		GroupName:     childGroup2.GroupName,
		Description:   childGroup2.Description,
		Extra:         childGroup2.Extra,
	})
	require.NoError(t, err)
	childGroup2.GroupId = createGroupResponse.GroupId
	childGroup2.ParentGroupId = parentGroup.GroupId
	childGroup2.GroupPath = childGroup2.ParentGroupId + "." + childGroup2.GroupId

	// modify child group2, move child group2 to child group1
	childGroup2.GroupName = "new test child2"
	childGroup2.Description = "new for test child2"
	childGroup2.Extra = map[string]string{
		"num": "16",
	}
	childGroup2.ParentGroupId = childGroup1.GroupId
	childGroup2.GroupPath = childGroup1.GroupPath + "." + childGroup2.GroupId
	_, err = imClient.ModifyGroup(ctx, &pb.ModifyGroupRequest{
		GroupId:       childGroup2.GroupId,
		ParentGroupId: childGroup2.ParentGroupId,
		GroupName:     childGroup2.GroupName,
		Description:   childGroup2.Description,
		Extra:         childGroup2.Extra,
	})
	require.NoError(t, err)

	// list child group2, use parent group id
	listGroupsResponse, err = imClient.ListGroups(ctx, &pb.ListGroupsRequest{
		ParentGroupId: []string{childGroup1.GroupId},
		Status:        []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listGroupsResponse.Total, 1)
	isGroupEqual(t, childGroup2, listGroupsResponse.GroupSet[0], constants.StatusActive)

	// delete parent group, has child group, can not delete
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{parentGroup.GroupId},
	})
	require.Error(t, err)

	// delete child group1, has child group, can not delete
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{childGroup1.GroupId},
	})
	require.Error(t, err)

	// delete child group2
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{childGroup2.GroupId},
	})
	require.NoError(t, err)
	listGroupsResponse, err = imClient.ListGroups(ctx, &pb.ListGroupsRequest{
		GroupId: []string{childGroup2.GroupId},
		Status:  []string{constants.StatusDeleted},
	})
	require.NoError(t, err)
	require.EqualValues(t, listGroupsResponse.Total, 1)
	isGroupEqual(t, childGroup2, listGroupsResponse.GroupSet[0], constants.StatusDeleted)

	// delete child group1
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{childGroup1.GroupId},
	})
	require.NoError(t, err)

	// delete parent group
	_, err = imClient.DeleteGroups(ctx, &pb.DeleteGroupsRequest{
		GroupId: []string{parentGroup.GroupId},
	})
	require.NoError(t, err)
}
