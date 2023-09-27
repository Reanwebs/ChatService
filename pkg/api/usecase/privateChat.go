package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"chat/pkg/api/repository"
	"log"
	"time"

	"gorm.io/gorm"
)

type PrivateChatUsecase struct {
	PrivateChatRepo repository.PrivateChatRepoMethods
}

func NewPrivateChatUsecase(repo repository.PrivateChatRepoMethods) PrivateChatUsecase {
	return PrivateChatUsecase{
		PrivateChatRepo: repo,
	}
}

type PrivateChatUsecaseMethods interface {
	StartChat(models.PrivateChat) error
	PrivateChatList(models.GetChat) ([]models.PrivateChat, error)
	CreatePrivateChatHistory(models.PrivateChatHistory) error
	RetrivePrivateChatHistory(string, string) ([]models.PrivateChatHistory, error)
	RetriveRecievedChatHistory(string, string) ([]models.PrivateChatHistory, error)
}

func (r PrivateChatUsecase) StartChat(input models.PrivateChat) error {
	entity := domain.PrivateChat{
		Model:             gorm.Model{},
		UserID:            input.UserID,
		UserName:          input.UserName,
		RecipientID:       input.RecipientID,
		RecipientName:     input.RecipientName,
		RecipientAvatarID: input.RecipientAvatarID,
		NewRecipient:      true,
		StartAt:           time.Now(),
		LastSeen:          time.Now(),
	}
	if err := r.PrivateChatRepo.CreatePrivateChat(entity); err != nil {
		return err
	}
	return nil
}

func (r PrivateChatUsecase) PrivateChatList(input models.GetChat) ([]models.PrivateChat, error) {
	existingChatList, newChatList, err := r.PrivateChatRepo.GetChatList(input.UserID)
	if err != nil {
		return nil, err
	}
	for i := range newChatList {
		if newChatList[i].RecipientID == input.UserID {
			newChatList[i].RecipientAvatarID = ""
			newChatList[i].UserID, newChatList[i].RecipientID = newChatList[i].RecipientID, newChatList[i].UserID
			newChatList[i].UserName, newChatList[i].RecipientName = newChatList[i].RecipientName, newChatList[i].UserName
		}
	}
	combinedChatList := append(existingChatList, newChatList...)
	return combinedChatList, nil
}

func (r PrivateChatUsecase) CreatePrivateChatHistory(chat models.PrivateChatHistory) error {
	input := domain.PrivateChatHistory{
		UserName:    chat.UserName,
		UserID:      chat.UserID,
		RecipientID: chat.RecipientID,
		Text:        chat.Text,
		Status:      chat.Status,
		Time:        time.Now(),
	}

	if err := r.PrivateChatRepo.AddPrivateChatHistory(input); err != nil {
		return err
	}
	return nil
}

func (r PrivateChatUsecase) RetrivePrivateChatHistory(userID string, recipientID string) ([]models.PrivateChatHistory, error) {
	response := []models.PrivateChatHistory{}
	DataLimitDays := time.Now().AddDate(0, 0, -2)
	result, err := r.PrivateChatRepo.GetPrivateChatHistory(userID, recipientID, DataLimitDays)
	if err != nil {
		log.Println("PrivateChatHistoryRepo", err)
		return nil, err
	}

	for _, chat := range result {
		response = append(response, models.PrivateChatHistory{
			UserName:    chat.UserName,
			UserID:      chat.UserID,
			RecipientID: chat.RecipientID,
			Text:        chat.Text,
			Status:      chat.Status,
			Time:        chat.Time,
		})
	}

	return response, nil
}

func (r PrivateChatUsecase) RetriveRecievedChatHistory(userID string, recipientID string) ([]models.PrivateChatHistory, error) {
	response := []models.PrivateChatHistory{}
	DataLimitDays := time.Now().AddDate(0, 0, -2)
	result, err := r.PrivateChatRepo.GetRecievedChatHistory(userID, recipientID, DataLimitDays)
	if err != nil {
		log.Println("PrivateChatHistoryRepo", err)
		return nil, err
	}

	for _, chat := range result {
		response = append(response, models.PrivateChatHistory{
			UserName:    chat.UserName,
			UserID:      chat.UserID,
			RecipientID: chat.RecipientID,
			Text:        chat.Text,
			Status:      chat.Status,
			Time:        chat.Time,
		})
	}

	return response, nil
}
