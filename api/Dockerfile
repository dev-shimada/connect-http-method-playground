FROM --platform=$BUILDPLATFORM golang:1.24.3-bookworm AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod /workspace
COPY go.sum /workspace
RUN go mod download
COPY ent /workspace/ent
COPY gen /workspace/gen
COPY main.go /workspace

RUN  <<EOF
# grpc_health_probe
go install github.com/grpc-ecosystem/grpc-health-probe@latest
CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/app main.go
EOF

FROM --platform=$BUILDPLATFORM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=build /workspace/bin/app /app/app
COPY --from=build /go/bin/grpc-health-probe /bin/grpc-health-probe
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s CMD ["/bin/grpc-health-probe", "-addr=:8081"]
EXPOSE 8081
ENTRYPOINT [ "/app/app" ]
