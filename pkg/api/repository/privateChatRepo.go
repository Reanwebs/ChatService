package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type PrivateChatRepo struct {
	DB mongo.Client
}

func NewPrivateChatRepo(dbClient mongo.Client) PrivateChatRepoMethods {
	return PrivateChatRepo{
		DB: dbClient,
	}
}

type PrivateChatRepoMethods interface {
	CreatePrivateChat()
}

func (r PrivateChatRepo) CreatePrivateChat() {
	fmt.Println("injected")
}
