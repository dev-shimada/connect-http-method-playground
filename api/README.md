# gRPC message sample
## Project Structure

- `proto/message/`: Contains the protobuf definitions and generated code.
- `buf.yaml`: Configuration for Buf, a tool for managing Protobuf files.
- `buf.gen.yaml`: Configuration for generating code from Protobuf files.

## Setup

### Prerequisites

- Docker

### Generating Protobuf Code

To generate the protobuf code using Buf, run:
```bash
buf generate
```

## grpcurl
```console
# grpcurl -plaintext localhost:8081 list
api.v1.ApiService
```
```console
# grpcurl -plaintext localhost:8081 list api.v1.ApiService
api.v1.ApiService.Get
api.v1.ApiService.PingPong
api.v1.ApiService.Post
```
```console
# grpcurl -plaintext -d '{"service": "api.v1.ApiService"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"user_id": "test_id", "text": "hello!"}' localhost:8081 api.v1.ApiService.PingPong
{
  "userId": "test_id",
  "text": "hello!"
}
```
```console
# uuidgen
210af8e1-9aee-4d2a-8751-e7d2cd9d39d3
# grpcurl -plaintext -d '{"user_id": "210af8e1-9aee-4d2a-8751-e7d2cd9d39d3", "text": "hello!"}' localhost:8081 api.v1.ApiService.Post
{
  "id": "01975d1b-f11a-77c0-a979-9f7dc73436e9"
}
```
```console
# grpcurl -plaintext -d '{"id": "01975d1b-f11a-77c0-a979-9f7dc73436e9"}' localhost:8081 api.v1.ApiService.Get
{
  "userId": "210af8e1-9aee-4d2a-8751-e7d2cd9d39d3",
  "text": "hello!"
}
```

## curl
```console
# curl --json '{"service": "api.v1.ApiService"}' localhost:8081/grpc.health.v1.Health/Check
{"status":"SERVING_STATUS_SERVING"}
```
```console
# curl --json '{"user_id": "test_id", "text": "hello!"}' localhost:8081/api.v1.ApiService/PingPong
{"userId":"test_id","text":"hello!"}
```
```console
# uuidgen
a4ddbcb1-7a8d-4f25-8a94-ea5ca7fe0b4b
# curl --json '{"user_id": "a4ddbcb1-7a8d-4f25-8a94-ea5ca7fe0b4b", "text": "hello! from curl"}' localhost:8081/api.v1.ApiService/Post
{"id":"01975d21-0d2a-7ee8-85de-ba5764e23a1d"}
```
```console
# curl --json '{"id": "01975d21-0d2a-7ee8-85de-ba5764e23a1d"}' localhost:8081/api.v1.ApiService/Get
{"userId":"a4ddbcb1-7a8d-4f25-8a94-ea5ca7fe0b4b","text":"hello! from curl"}
```

## http/1.1
```console
# curl "localhost:8081/pingpong?user_id=test_id&text=hello!"
{"userId":"test_id","text":"hello!"}
```
