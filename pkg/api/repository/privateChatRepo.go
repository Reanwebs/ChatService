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
	GetChatList(string) ([]models.PrivateChat, error)
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
		if result := r.DB.Model(&existingRecord).Update("LastSeen", time.Now()); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}

	return nil
}

func (r PrivateChatRepo) GetChatList(userID string) ([]models.PrivateChat, error) {
	var chatList []models.PrivateChat
	if err := r.DB.Where("user_id = ?", userID).Find(&chatList).Error; err != nil {
		return nil, err
	}
	return chatList, nil
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
	err := r.DB.Where("user_id = ? AND recipient_id = ?", userID, recipientID).Where("time >= ?", dateLimit).
		Find(&chats).Error

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
