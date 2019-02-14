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

package im

import (
	"context"

	"kubesphere.io/im/pkg/pb"
	"kubesphere.io/im/pkg/service/im/resource"
	"kubesphere.io/im/pkg/version"
)

func (p *Server) GetVersion(ctx context.Context, req *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	reply := &pb.GetVersionResponse{Version: version.GetVersionString()}
	return reply, nil
}

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	return resource.CreateGroup(ctx, req)
}

func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	return resource.DeleteGroups(ctx, req)
}

func (p *Server) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	group, err := resource.GetGroup(ctx, req.GroupId)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetGroupResponse{
			Group: group.ToPB(),
		}, nil
	}
}

func (p *Server) GetGroupWithUser(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupWithUserResponse, error) {
	groupWithUser, err := resource.GetGroupWithUser(ctx, req.GroupId)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetGroupWithUserResponse{
			Group: groupWithUser.ToPB(),
		}, nil
	}
}

func (p *Server) ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	return resource.ListGroups(ctx, req)
}

func (p *Server) ListGroupsWithUser(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsWithUserResponse, error) {
	return resource.ListGroupsWithUser(ctx, req)
}

func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	return resource.ModifyGroup(ctx, req)
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return resource.CreateUser(ctx, req)
}

func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	return resource.DeleteUsers(ctx, req)
}

func (p *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := resource.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetUserResponse{
			User: user.ToPB(),
		}, nil
	}
}

func (p *Server) GetUserWithGroup(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserWithGroupResponse, error) {
	userWithGroup, err := resource.GetUserWithGroup(ctx, req.UserId)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetUserWithGroupResponse{
			User: userWithGroup.ToPB(),
		}, nil
	}
}

func (p *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return resource.ListUsers(ctx, req)
}

func (p *Server) ListUsersWithGroup(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersWithGroupResponse, error) {
	return resource.ListUsersWithGroup(ctx, req)
}

func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	return resource.ModifyUser(ctx, req)
}

func (p *Server) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	return resource.JoinGroup(ctx, req)
}

func (p *Server) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	return resource.LeaveGroup(ctx, req)
}

func (p *Server) ComparePassword(ctx context.Context, req *pb.ComparePasswordRequest) (*pb.ComparePasswordResponse, error) {
	return resource.ComparePassword(ctx, req)
}

func (p *Server) ModifyPassword(ctx context.Context, req *pb.ModifyPasswordRequest) (*pb.ModifyPasswordResponse, error) {
	return resource.ModifyPassword(ctx, req)
}
