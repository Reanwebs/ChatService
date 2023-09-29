package repository

import (
	"chat/pkg/api/domain"
	"log"
	"time"

	"gorm.io/gorm"
)

type GroupChatRepo struct {
	DB *gorm.DB
}

func NewGroupChatRepo(dbClient *gorm.DB) GroupChatRepoMethods {
	return GroupChatRepo{
		DB: dbClient,
	}
}

type GroupChatRepoMethods interface {
	CreateGroupChat(domain.GroupChat) error
	GetGroupList(domain.GroupChat) ([]domain.GroupChat, error)
	AddGroupChatHistory(domain.GroupChatHistory) error
	GroupChatHistory(domain.GroupChatHistory, time.Time) ([]domain.GroupChatHistory, error)
}

func (r GroupChatRepo) CreateGroupChat(input domain.GroupChat) error {
	var groupChat domain.GroupChat
	if result := r.DB.Where("user_id = ? AND group_id = ?", input.UserID, input.GroupID).First(&groupChat); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

			if result := r.DB.Create(&input); result.Error != nil {
				log.Println(result.Error)
				return result.Error
			}
			return nil
		}
		log.Println(result.Error)
		return result.Error
	} else {
		updateData := domain.GroupChat{
			Model:         gorm.Model{},
			UserID:        input.UserID,
			UserName:      input.UserName,
			GroupID:       input.GroupID,
			GroupName:     input.GroupName,
			GroupAvatarID: input.GroupAvatarID,
			Permission:    input.Permission,
			StartAt:       time.Time{},
			LastSeen:      time.Now(),
		}
		if result := r.DB.Model(&groupChat).Updates(updateData); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
	}

	return nil
}

func (r GroupChatRepo) GetGroupList(input domain.GroupChat) ([]domain.GroupChat, error) {
	var groupList []domain.GroupChat
	if err := r.DB.Where("user_id = ?", input.UserID).Find(&groupList).Error; err != nil {
		return nil, err
	}
	return groupList, nil
}
func (r GroupChatRepo) AddGroupChatHistory(chatHistory domain.GroupChatHistory) error {
	if err := r.DB.Create(&chatHistory).Error; err != nil {
		return err
	}
	return nil
}

func (r GroupChatRepo) GroupChatHistory(input domain.GroupChatHistory, dateLimit time.Time) ([]domain.GroupChatHistory, error) {
	var chatHistory []domain.GroupChatHistory
	err := r.DB.Where("group_id = ?", input.GroupID).Where("time >= ?", dateLimit).Find(&chatHistory).Error
	if err != nil {
		return nil, err
	}
	return chatHistory, nil
}
