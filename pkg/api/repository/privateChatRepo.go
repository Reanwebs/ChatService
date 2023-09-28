package repository

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"fmt"
	"log"
	"time"

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
	GetChatList(string) ([]models.PrivateChat, []models.PrivateChat, error)
	AddPrivateChatHistory(domain.PrivateChatHistory) error
	GetPrivateChatHistory(string, string, time.Time) ([]domain.PrivateChatHistory, error)
	GetRecievedChatHistory(string, string, time.Time) ([]domain.PrivateChatHistory, error)
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
	} else {
		updateData := domain.PrivateChat{
			UserID:            input.UserID,
			UserName:          input.UserName,
			RecipientID:       input.RecipientID,
			RecipientName:     input.RecipientName,
			RecipientAvatarID: input.RecipientAvatarID,
			NewRecipient:      input.NewRecipient,
			LastSeen:          time.Now(),
		}
		if result := r.DB.Model(&existingRecord).Updates(updateData); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}
	return nil
}

func (r PrivateChatRepo) GetChatList(userID string) ([]models.PrivateChat, []models.PrivateChat, error) {
	var existingChatList []models.PrivateChat
	var newChatList []models.PrivateChat
	if err := r.DB.Where("user_id = ?", userID).Find(&existingChatList).Error; err != nil {
		return nil, nil, err
	}
	if err := r.DB.Where("recipient_id = ? AND new_recipient = ?", userID, true).Find(&newChatList).Error; err != nil {
		return nil, nil, err
	}
	if result := r.DB.Model(&models.PrivateChat{}).Where("recipient_id = ? AND new_recipient = ?", userID, true).
		Updates(map[string]interface{}{"new_recipient": false}); result.Error != nil {
		return nil, nil, result.Error
	}

	return existingChatList, newChatList, nil
}

func (r PrivateChatRepo) AddPrivateChatHistory(chat domain.PrivateChatHistory) error {
	result := r.DB.Create(&chat)

	if result.Error != nil {

		fmt.Println(result.Error)
	}
	return nil
}

func (r PrivateChatRepo) GetPrivateChatHistory(userID string, recipientID string, dateLimit time.Time) ([]domain.PrivateChatHistory, error) {
	var chats []domain.PrivateChatHistory
	err := r.DB.Where("user_id = ? AND recipient_id = ?", userID, recipientID).Where("time >= ?", dateLimit).Find(&chats).Error
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (r PrivateChatRepo) GetRecievedChatHistory(userID string, recipientID string, dateLimit time.Time) ([]domain.PrivateChatHistory, error) {
	var chats []domain.PrivateChatHistory
	err := r.DB.Where("user_id = ? AND recipient_id = ?", recipientID, userID).Where("time >= ?", dateLimit).
		Find(&chats).Error

	if err != nil {
		return nil, err
	}

	return chats, nil
}
