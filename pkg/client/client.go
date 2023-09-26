package client

import (
	"chat/pb/client"
	authclient "chat/pkg/api/delivery/authClient"
	"chat/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitializeClient(c server.Config) (authclient.AutharizationClientMethods, error) {
	clientCon, err := grpc.Dial(c.AuthClientURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return authclient.NewAuthClient(client.NewAutharizationClient(clientCon)), nil
}
