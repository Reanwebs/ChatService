package repository

import (
	"chat/pkg/api/delivery/models"
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
	GetChatList(string) ([]models.PrivateChat, error)
}

func (r PrivateChatRepo) GetChatList(userID string) ([]models.PrivateChat, error) {
	var chatList []models.PrivateChat
	if err := r.DB.Where("user_id = ?", userID).Find(&chatList).Error; err != nil {
		return nil, err
	}
	return chatList, nil
}

func (r PrivateChatRepo) CreatePrivateChat(input domain.PrivateChat) error {
	var existingRecord domain.PrivateChat

	if result := r.DB.Where("user_id = ? AND recipient_id = ?", input.UserID, input.RecipientID).First(&existingRecord); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if result := r.DB.Create(&input); result.Error != nil {
				log.Println(result.Error)
				return result.Error
			}
		} else {
			log.Println(result.Error)
			return result.Error
		}
	}

	return nil
}
