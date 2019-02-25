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

func isUserEqual(t *testing.T, oldUser, newUser *pb.User, status string) bool {
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.PhoneNumber, newUser.PhoneNumber)
	require.Equal(t, oldUser.Description, newUser.Description)
	require.Equal(t, oldUser.Extra["age"], newUser.Extra["age"])
	require.Equal(t, status, newUser.Status)
	return true
}

func TestUser(t *testing.T) {
	prepare(t)

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

	// modify password
	password = "newpassw0rd"
	_, err = imClient.ModifyPassword(ctx, &pb.ModifyPasswordRequest{
		UserId:   user.UserId,
		Password: password,
	})
	require.NoError(t, err)

	// compare password
	comparePasswordResponse, err := imClient.ComparePassword(ctx, &pb.ComparePasswordRequest{
		UserId:   user.UserId,
		Password: password,
	})
	require.NoError(t, err)
	require.EqualValues(t, comparePasswordResponse.Ok, true)

	// get user
	getUserResponse, err := imClient.GetUser(ctx, &pb.GetUserRequest{
		UserId: user.UserId,
	})
	require.NoError(t, err)
	isUserEqual(t, user, getUserResponse.User, constants.StatusActive)

	// list user, use email
	listUsersResponse, err := imClient.ListUsers(ctx, &pb.ListUsersRequest{
		Email:  []string{user.Email},
		Status: []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listUsersResponse.Total, 1)
	isUserEqual(t, user, listUsersResponse.UserSet[0], constants.StatusActive)

	// list user, use user id
	listUsersResponse, err = imClient.ListUsers(ctx, &pb.ListUsersRequest{
		UserId: []string{user.UserId},
		Status: []string{constants.StatusActive},
	})
	require.NoError(t, err)
	require.EqualValues(t, listUsersResponse.Total, 1)
	isUserEqual(t, user, listUsersResponse.UserSet[0], constants.StatusActive)

	// modify user
	user.Username = "new test"
	user.Email = "new_test@op.com"
	user.PhoneNumber = "11111111111"
	user.Description = "for new test"
	user.Extra = map[string]string{
		"age": "21",
	}
	_, err = imClient.ModifyUser(ctx, &pb.ModifyUserRequest{
		UserId:      user.UserId,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Description: user.Description,
		Extra:       user.Extra,
	})
	require.NoError(t, err)
	getUserResponse, err = imClient.GetUser(ctx, &pb.GetUserRequest{
		UserId: user.UserId,
	})
	require.NoError(t, err)
	isUserEqual(t, user, getUserResponse.User, constants.StatusActive)

	// delete user
	_, err = imClient.DeleteUsers(ctx, &pb.DeleteUsersRequest{
		UserId: []string{user.UserId},
	})
	require.NoError(t, err)
	listUsersResponse, err = imClient.ListUsers(ctx, &pb.ListUsersRequest{
		UserId: []string{user.UserId},
		Status: []string{constants.StatusDeleted},
	})
	require.NoError(t, err)
	require.EqualValues(t, listUsersResponse.Total, 1)
	isUserEqual(t, user, listUsersResponse.UserSet[0], constants.StatusDeleted)

}
