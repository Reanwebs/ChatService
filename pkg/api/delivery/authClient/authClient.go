package authclient

import (
	"chat/pb/client"
	"context"
	"log"
	"time"
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
	GetUserDetails(context.Context, *client.GetUserDetailsRequest) (*client.GetUserDetailsResponse, error)
	GetOnlineStatus(context.Context, *client.GetOnlineStatusRequest) (*client.GetOnlineStatusResponse, error)
	UserGroupPermission(context.Context, *client.UserGroupPermissionRequest) (*client.UserGroupPermissionResponse, error)
}

func (c AutharizationClient) GetOnlineStatus(ctx context.Context, req *client.GetOnlineStatusRequest) (*client.GetOnlineStatusResponse, error) {
	return nil, nil
}

func (c AutharizationClient) HealthCheck(ctx context.Context, req *client.Request) (*client.Response, error) {
	return nil, nil
}

func (c AutharizationClient) GetUserDetails(ctx context.Context, req *client.GetUserDetailsRequest) (*client.GetUserDetailsResponse, error) {
	var res *client.GetUserDetailsResponse
	var err error
	for retry := 0; retry < 3; retry++ {
		res, err = c.Client.GetUserDetails(ctx, &client.GetUserDetailsRequest{
			UserID: req.UserID,
		})
		if err == nil || ctx.Err() != nil {
			break
		}
		log.Println("retrying-client")
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c AutharizationClient) UserGroupPermission(ctx context.Context, req *client.UserGroupPermissionRequest) (*client.UserGroupPermissionResponse, error) {
	var res *client.UserGroupPermissionResponse
	var err error
	for retry := 0; retry < 3; retry++ {
		res, err = c.Client.UserGroupPermission(ctx, &client.UserGroupPermissionRequest{
			UserID:  req.UserID,
			GroupID: req.GroupID,
		})
		if err == nil || ctx.Err() != nil {
			break
		}
		log.Println("retrying-client")
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
