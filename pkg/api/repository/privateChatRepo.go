package repository

import (
	"chat/pkg/api/domain"
	"log"

	"gorm.io/gorm"
)

type PrivateChatRepo struct {
	DB *gorm.DB
}

func NewPrivateChatRepo(dbClient *gorm.DB) PrivateChatRepoMethods {
	return PrivateChatRepo{
		DB: dbClient,
	}
}

type PrivateChatRepoMethods interface {
	CreatePrivateChat(domain.PrivateChat) error
}

func (r PrivateChatRepo) CreatePrivateChat(input domain.PrivateChat) error {
	if result := r.DB.Create(&input); result.Error != nil {
		log.Println(result.Error)
	}
	return nil
}
