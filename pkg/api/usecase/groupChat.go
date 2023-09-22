package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"chat/pkg/api/repository"
	"errors"
	"log"
	"time"
)

type GroupChatUsecase struct {
	GroupChatRepo repository.GroupChatRepoMethods
}

func NewGroupChatUsecase(repo repository.GroupChatRepoMethods) GroupChatUsecase {
	return GroupChatUsecase{
		GroupChatRepo: repo,
	}
}

type GroupChatUsecaseMethods interface {
	GroupChatStart(models.GroupChat) error
	GetGroupList(models.GetGroupChat) ([]models.GroupChat, error)
}

func (u GroupChatUsecase) GroupChatStart(input models.GroupChat) error {
	entity := domain.GroupChat{
		UserID:   input.UserID,
		GroupID:  input.GroupID,
		StartAt:  time.Now(),
		LastSeen: time.Time{},
	}
	if err := u.GroupChatRepo.CreateGroupChat(entity); err != nil {
		return errors.Join(err, errors.New("error in starting chat"))
	}
	return nil
}

func (u GroupChatUsecase) GetGroupList(input models.GetGroupChat) ([]models.GroupChat, error) {
	entity := domain.GroupChat{
		UserID: input.UserID,
	}
	response, err := u.GroupChatRepo.GetGroupList(entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return response, nil
}
