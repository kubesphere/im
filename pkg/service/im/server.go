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
	"os"

	"github.com/google/gops/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/config"
	"kubesphere.io/im/pkg/global"
	"kubesphere.io/im/pkg/manager"
	"kubesphere.io/im/pkg/pb"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	global.SetGlobal(cfg)
	s := new(Server)
	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true,
	}); err != nil {
		logger.Criticalf(nil, "failed to start gops agent")
	}
	if cfg.TlsEnabled {
		creds, err := credentials.NewServerTLSFromFile(cfg.TlsCertFile, cfg.TlsKeyFile)
		if err != nil {
			logger.Criticalf(nil, "Constructs TLS credentials failed: %+v", err)
			os.Exit(1)
		}
		manager.NewGrpcServer(cfg.Host, cfg.Port).
			Serve(func(server *grpc.Server) {
				pb.RegisterIdentityManagerServer(server, s)
				grpc.Creds(creds)
			})
	} else {
		manager.NewGrpcServer(cfg.Host, cfg.Port).
			Serve(func(server *grpc.Server) {
				pb.RegisterIdentityManagerServer(server, s)
			})
	}
}
