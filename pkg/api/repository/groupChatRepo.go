package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type GroupChatRepo struct {
	DB mongo.Client
}

func NewGroupChatRepo(dbClient mongo.Client) GroupChatRepoMethods {
	return GroupChatRepo{
		DB: dbClient,
	}
}

type GroupChatRepoMethods interface {
	CreateGroupChat()
}

func (r GroupChatRepo) CreateGroupChat() {
	fmt.Println("injected")
}
