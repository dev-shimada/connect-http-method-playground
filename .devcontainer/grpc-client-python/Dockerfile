FROM --platform=$BUILDPLATFORM python:3.14-trixie AS vscode

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ENV ENV PYTHONPATH=/workspace/grpc-client-python/gen

WORKDIR /workspace/grpc-client-python

RUN  <<EOF
apt-get update && apt-get install -y vim git unzip less bash-completion jq sqlite3 uuid-runtime npm protobuf-compiler
npm i -g @anthropic-ai/claude-code
# buf cli
BIN="/usr/local/bin"
VERSION="1.35.1"
curl -sSL "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" -o "${BIN}/buf"
chmod +x "${BIN}/buf"
# bash completions
echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc
echo "source <(buf completion bash)" >> ~/.bashrc
EOF

COPY pyproject.toml /workspace/grpc-client-python/

RUN <<EOF
pip isntall uv
uv pip install --system -r pyproject.toml
EOF

COPY . /workspace/grpc-client-python/
