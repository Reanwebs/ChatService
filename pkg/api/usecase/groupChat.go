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
	GetGroupChatHistory(models.GroupChat) ([]models.GroupChatHistory, error)
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
	var convertedResponse []models.GroupChat
	for _, group := range response {
		convertedResponse = append(convertedResponse, models.GroupChat{
			UserID:   group.UserID,
			GroupID:  group.GroupID,
			StartAt:  group.StartAt,
			LastSeen: group.LastSeen,
		})
	}
	return convertedResponse, nil
}

func (u GroupChatUsecase) GetGroupChatHistory(input models.GroupChat) ([]models.GroupChatHistory, error) {
	entity := domain.GroupChat{
		UserID:  input.UserID,
		GroupID: input.GroupID,
	}
	DataLimitDays := time.Now().AddDate(0, 0, -2)
	response, err := u.GroupChatRepo.GroupChatHistory(entity, DataLimitDays)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var convertedResponse []models.GroupChatHistory
	for _, groupChat := range response {
		convertedResponse = append(convertedResponse, models.GroupChatHistory{
			GroupID: groupChat.GroupID,
			Text:    groupChat.Text,
			Status:  groupChat.Status,
			Time:    time.Time{},
		})
	}
	return convertedResponse, nil
}
