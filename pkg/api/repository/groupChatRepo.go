package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type GroupChatRepo struct {
	DB *gorm.DB
}

func NewGroupChatRepo(dbClient *gorm.DB) GroupChatRepoMethods {
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
