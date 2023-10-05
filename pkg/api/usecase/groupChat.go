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
	GetGroupList(string) ([]models.GroupChat, error)
	AddGroupChatHistory(models.GroupChatHistory) error
	GetGroupChatHistory(models.GroupChatHistory) ([]models.GroupChatHistory, error)
}

func (u GroupChatUsecase) GroupChatStart(input models.GroupChat) error {
	entity := domain.GroupChat{
		UserID:        input.UserID,
		UserName:      input.UserName,
		GroupID:       input.GroupID,
		GroupName:     input.GroupName,
		GroupAvatarID: input.GroupAvatarID,
		Permission:    input.Permission,
		StartAt:       time.Now(),
		LastSeen:      time.Time{},
	}
	if err := u.GroupChatRepo.CreateGroupChat(entity); err != nil {
		return errors.Join(err, errors.New("error in starting chat"))
	}
	return nil
}

func (u GroupChatUsecase) GetGroupList(input string) ([]models.GroupChat, error) {
	entity := domain.GroupChat{
		UserID: input,
	}
	response, err := u.GroupChatRepo.GetGroupList(entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var convertedResponse []models.GroupChat
	for _, group := range response {
		convertedResponse = append(convertedResponse, models.GroupChat{
			UserID:        group.UserID,
			UserName:      group.UserName,
			GroupID:       group.GroupID,
			GroupName:     group.GroupName,
			GroupAvatarID: group.GroupAvatarID,
			Permission:    group.Permission,
			StartAt:       group.StartAt,
			LastSeen:      group.LastSeen,
		})
	}
	return convertedResponse, nil
}

func (u GroupChatUsecase) AddGroupChatHistory(input models.GroupChatHistory) error {
	entity := domain.GroupChatHistory{
		UserID:    input.UserID,
		UserName:  input.UserName,
		GroupID:   input.GroupID,
		GroupName: input.GroupName,
		Text:      input.Text,
		Status:    input.Status,
		Time:      time.Now(),
	}
	if err := u.GroupChatRepo.AddGroupChatHistory(entity); err != nil {
		return err
	}
	return nil
}

func (u GroupChatUsecase) GetGroupChatHistory(input models.GroupChatHistory) ([]models.GroupChatHistory, error) {
	entity := domain.GroupChatHistory{
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
			UserID:    groupChat.UserID,
			UserName:  groupChat.UserName,
			GroupID:   groupChat.GroupID,
			GroupName: groupChat.GroupName,
			Text:      groupChat.Text,
			Status:    groupChat.Status,
			Time:      time.Time{},
		})
	}
	return convertedResponse, nil
}
