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

package constants

const (
	ColumnUserId         = "user_id"
	ColumnGroupId        = "group_id"
	ColumnCreateTime     = "create_time"
	ColumnUpdateTime     = "update_time"
	ColumnStatusTime     = "status_time"
	ColumnStatus         = "status"
	ColumnPassword       = "password"
	ColumnEmail          = "email"
	ColumnPhoneNumber    = "phone_number"
	ColumnGroupPath      = "group_path"
	ColumnUsername       = "username"
	ColumnGroupName      = "group_name"
	ColumnParentGroupId  = "parent_group_id"
	ColumnGroupPathLevel = "group_path_level"
	ColumnDescription    = "description"
	ColumnExtra          = "extra"
)

const (
	TableUserGroupBinding = "user_group_binding"
	TableUser             = "user"
	TableGroup            = "group"
)

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableUser: {
		ColumnUserId, ColumnEmail, ColumnPhoneNumber, ColumnStatus,
	},
	TableGroup: {
		ColumnGroupId, ColumnGroupPath, ColumnStatus,
	},
}

var SearchWordColumnTable = []string{
	TableUser,
	TableGroup,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	TableUser: {
		ColumnUserId, ColumnUsername, ColumnEmail, ColumnPhoneNumber, ColumnStatus,
	},
	TableGroup: {
		ColumnGroupId, ColumnGroupName, ColumnGroupPath, ColumnStatus,
	},
}
