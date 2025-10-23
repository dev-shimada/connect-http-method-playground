FROM --platform=$BUILDPLATFORM golang:1.24.3-bookworm AS vscode

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ENV LESSCHARSET utf-8

WORKDIR /workspace

RUN  <<EOF
apt-get update && apt-get install -y vim git unzip less bash-completion jq sqlite3 uuid-runtime
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/cweill/gotests/gotests@latest
go install github.com/goreleaser/goreleaser/v2@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
# buf cli
BIN="/usr/local/bin"
VERSION="1.35.1"
curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" -o "${BIN}/buf"
chmod +x "${BIN}/buf"
# protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# connect-go
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
# grpc_health_probe
go install github.com/grpc-ecosystem/grpc-health-probe@latest
echo export PATH="$PATH:$(go env GOPATH)/bin" >> ~/.bashrc
# bash completions
echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc
echo "source <(buf completion bash)" >> ~/.bashrc
EOF

COPY . /workspace
