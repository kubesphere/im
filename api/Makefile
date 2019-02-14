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

PWD:=$(shell pwd)

GOPATH:=$(shell go env GOPATH)

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(PROTOC_INC_PATH)

generate: $(PROTO_FILES) Makefile
	@rm -rf ../pkg/pb/*
	@mkdir -p ../pkg/pb

	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg
	# format generated files
	goimports -l -w -e -local=openpitrix ../pkg/pb

clean:
	rm -rf $(GOPATH)/src/kubesphere.io/im/pkg/pb
