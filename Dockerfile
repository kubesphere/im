# Copyright 2019 The KubeSphere Authors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.11-alpine3.7 as builder

# install tools
RUN apk add --no-cache git

WORKDIR /go/src/kubesphere.io/im
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir -p /kubesphere_bin
RUN go generate kubesphere.io/im/pkg/version && \
	GOBIN=/kubesphere_bin go install -ldflags '-w -s' -tags netgo kubesphere.io/im/cmd/...

FROM alpine:3.7
COPY --from=builder /kubesphere_bin/im /usr/local/bin/
CMD ["/usr/local/bin/im"]
