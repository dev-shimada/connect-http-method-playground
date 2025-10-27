package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/dev-shimada/connect-http-method-playground/grpc-client/gen/proto/api/v1"
	"github.com/dev-shimada/connect-http-method-playground/grpc-client/gen/proto/api/v1/apiv1connect"
	"github.com/google/uuid"
)

const (
	serverAddr string = "api:8081"
)

func main() {
	client := apiv1connect.NewApiServiceClient(http.DefaultClient, "http://"+serverAddr)
	ctx := context.Background()
	uuid := uuid.NewString()
	id, err := post(ctx, client, uuid)
	if err != nil {
		panic(err)
	}
	getRes, err := get(ctx, client, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get response: %+v", getRes)
}

func post(ctx context.Context, client apiv1connect.ApiServiceClient, uuid string) (string, error) {
	req := connect.NewRequest(&v1.PostRequest{
		UserId: uuid,
		Text:   "test",
	})
	res, err := client.Post(ctx, req)
	if err != nil {
		return "", err
	}
	println(res.Msg.Id)
	return res.Msg.Id, nil
}

func get(ctx context.Context, client apiv1connect.ApiServiceClient, id string) (*v1.GetResponse, error) {
	req := connect.NewRequest(&v1.GetRequest{
		Id: id,
	})
	res, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Msg, nil
}
