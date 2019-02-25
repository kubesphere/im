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
	"testing"

	"github.com/stretchr/testify/require"

	"kubesphere.io/im/pkg/client/im"
	"kubesphere.io/im/pkg/config"
	"kubesphere.io/im/pkg/global"
)

var imClient *im.Client

func prepare(t *testing.T) {
	cfg := config.Default()
	cfg.DB.Host = "127.0.0.1"
	cfg.Host = "127.0.0.1"
	cfg.DB.Port = 13306
	global.SetGlobal(cfg)
	var err error
	imClient, err = im.NewClient()
	require.NoError(t, err)
}
