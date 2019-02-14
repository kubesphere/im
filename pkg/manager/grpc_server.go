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

package manager

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/version"
)

type GrpcServer struct {
	ServiceName string
	Port        int
}

type RegisterCallback func(*grpc.Server)

func NewGrpcServer(serviceName string, port int) *GrpcServer {
	return &GrpcServer{
		ServiceName: serviceName,
		Port:        port,
	}
}

func (g *GrpcServer) Serve(callback RegisterCallback, opt ...grpc.ServerOption) {
	logger.Infof(nil, "Release version: %s", version.GetVersionString())
	logger.Infof(nil, "Service [%s] start listen at port [%d]", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Criticalf(nil, "failed to listen: %+v", err)
	}

	builtinOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			g.unaryServerLogInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Criticalf(nil, "GRPC server recovery with error: %+v", p)
					logger.Criticalf(nil, string(debug.Stack()))
					if e, ok := p.(error); ok {
						return e
					}
					return status.Errorf(codes.Internal, "panic")
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Criticalf(nil, "GRPC server recovery with error: %+v", p)
					logger.Criticalf(nil, string(debug.Stack()))
					if e, ok := p.(error); ok {
						return e
					}
					return status.Errorf(codes.Internal, "panic")
				}),
			),
		),
	}

	grpcServer := grpc.NewServer(append(opt, builtinOptions...)...)
	reflection.Register(grpcServer)
	callback(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Criticalf(nil, "%+v", err)
	}
}

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		OrigName: true,
	}
)

func (g *GrpcServer) unaryServerLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error

		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]
		if p, ok := req.(proto.Message); ok {
			if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
				logger.Errorf(ctx, "Failed to marshal proto message to string [%s][%+v]", action, err)
			} else {
				logger.Infof(ctx, "Request received [%s] [%s]", action, content)
			}
		}
		start := time.Now()

		resp, err := handler(ctx, req)

		elapsed := time.Since(start)
		logger.Infof(ctx, "Handled request [%s] exec_time is [%s]", action, elapsed)
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				logger.Debugf(ctx, "Response is error: %s, %s", e.Code().String(), e.Message())
			}
		}
		return resp, err
	}
}
