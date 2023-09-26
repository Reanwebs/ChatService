package authclient

import (
	"chat/pb/client"
	"context"
)

type AutharizationClient struct {
	Client client.AutharizationClient
}

func NewAuthClient(client client.AutharizationClient) AutharizationClientMethods {
	return &AutharizationClient{
		Client: client,
	}
}

type AutharizationClientMethods interface {
	HealthCheck(context.Context, *client.Request) (*client.Response, error)
	ValidateUser(context.Context, *client.ValidateUserRequest) (*client.ValidateUserResponse, error)
	GetOnlineStatus(context.Context, *client.GetOnlineStatusRequest) (*client.GetOnlineStatusResponse, error)
}

func (c AutharizationClient) GetOnlineStatus(ctx context.Context, req *client.GetOnlineStatusRequest) (*client.GetOnlineStatusResponse, error) {
	return nil, nil
}

func (c AutharizationClient) ValidateUser(ctx context.Context, req *client.ValidateUserRequest) (*client.ValidateUserResponse, error) {
	return nil, nil
}

func (c AutharizationClient) HealthCheck(ctx context.Context, req *client.Request) (*client.Response, error) {
	return nil, nil
}
