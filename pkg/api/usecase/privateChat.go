package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"chat/pkg/api/repository"
	"log"
	"time"
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
	CreatePrivateChatHistory(string, string, models.PrivateChatHistory) error
	RetrivePrivateChatHistory(string, string) ([]models.PrivateChatHistory, error)
	RetriveRecievedChatHistory(string, string) ([]models.PrivateChatHistory, error)
}

func (r PrivateChatUsecase) StartChat(input models.PrivateChat) error {
	entity := domain.PrivateChat{}

	entity.UserID = input.UserID
	entity.RecipientID = input.RecipientID
	entity.StartAt = time.Now()
	entity.LastSeen = time.Now()

	if err := r.PrivateChatRepo.CreatePrivateChat(entity); err != nil {
		return err
	}

	return nil
}

func (r PrivateChatUsecase) PrivateChatList(input models.GetChat) ([]models.PrivateChat, error) {
	response := []models.PrivateChat{}
	response, err := r.PrivateChatRepo.GetChatList(input.UserID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r PrivateChatUsecase) CreatePrivateChatHistory(userID string, recipientID string, chat models.PrivateChatHistory) error {
	input := domain.PrivateChatHistory{
		UserID:      userID,
		RecipientID: recipientID,
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
			UserID:      chat.UserID,
			RecipientID: chat.RecipientID,
			Text:        chat.Text,
			Status:      chat.Status,
			Time:        chat.Time,
		})
	}

	return response, nil
}
