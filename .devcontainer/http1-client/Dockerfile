FROM --platform=$BUILDPLATFORM golang:1.25.3-trixie AS vscode

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

RUN  <<EOF
apt-get update && apt-get install -y vim git unzip less bash-completion jq uuid-runtime npm
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/cweill/gotests/gotests@latest
go install github.com/goreleaser/goreleaser/v2@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
npm i -g @anthropic-ai/claude-code
echo export PATH="$PATH:$(go env GOPATH)/bin" >> ~/.bashrc
# bash completions
echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc
EOF

COPY . /workspace
