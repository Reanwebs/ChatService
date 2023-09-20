package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"chat/pkg/api/repository"
	"log"
	"sync"
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
	RetrivePrivateChatHistory(string, string) ([]models.PrivateChatWithHistory, error)
}

func (r PrivateChatUsecase) StartChat(input models.PrivateChat) error {
	entity := domain.PrivateChat{}

	entity.UserID = input.UserID
	entity.RecipientID = input.RecipientID
	entity.StartAt = time.Now()

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
		Time:        chat.Time,
	}

	if err := r.PrivateChatRepo.AddPrivateChatHistory(input); err != nil {
		return err
	}
	return nil
}

func (r PrivateChatUsecase) RetrivePrivateChatHistory(userID string, recipientID string) ([]models.PrivateChatWithHistory, error) {
	response := []models.PrivateChatWithHistory{}
	var wg sync.WaitGroup
	_, err := r.PrivateChatRepo.GetPrivateChatHistory(userID, recipientID)
	if err != nil {
		log.Println("PrivateChatHistoryRepo", err)
		return nil, err
	}

	responseChan := make(chan models.PrivateChatWithHistory)

	// for _, chat := range result {
	// 	wg.Add(1)
	// 	go func(chat domain.PrivateChatWithHistory) {
	// 		defer wg.Done()
	// 		responseChan <- MapDomainToModel(chat)
	// 	}(chat)
	// }

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	for mappedChat := range responseChan {
		response = append(response, mappedChat)
	}

	return response, nil
}

func MapDomainToModel(Chat domain.PrivateChatWithHistory) models.PrivateChatWithHistory {
	return models.PrivateChatWithHistory{
		PrivateChat: models.PrivateChat{
			UserID:      Chat.PrivateChat.UserID,
			RecipientID: Chat.PrivateChat.RecipientID,
		},
		PrivateChatHistory: models.PrivateChatHistory{
			Text:   Chat.PrivateChatHistory.Text,
			Status: Chat.PrivateChatHistory.Status,
			Time:   Chat.PrivateChatHistory.Time,
		},
	}
}
