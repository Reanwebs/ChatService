package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"chat/pkg/api/repository"
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
	PrivateChatList(models.GetChat) ([]models.PrivateChat, error)
	StartChat(models.PrivateChat) error
}

func (r PrivateChatUsecase) PrivateChatList(input models.GetChat) ([]models.PrivateChat, error) {
	response := []models.PrivateChat{}
	response, err := r.PrivateChatRepo.GetChatList(input.UserID)
	if err != nil {
		return nil, err
	}

	return response, nil
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
