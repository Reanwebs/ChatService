package repository

import (
	"fmt"
)

type PrivateChatRepo struct {
	DB interface{}
}

func NewPrivateChatRepo(dbClient interface{}) PrivateChatRepoMethods {
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
