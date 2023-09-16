package repository

import (
	"fmt"
)

type GroupChatRepo struct {
	DB interface{}
}

func NewGroupChatRepo(dbClient interface{}) GroupChatRepoMethods {
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
