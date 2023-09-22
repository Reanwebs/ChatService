package repository

import (
	"chat/pkg/api/delivery/models"
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
	GetGroupList(domain.GroupChat) ([]models.GroupChat, error)
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
	}

	if result := r.DB.Model(&groupChat).Update("LastSeen", time.Now()); result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	return nil
}

func (r GroupChatRepo) GetGroupList(input domain.GroupChat) ([]models.GroupChat, error) {
	var groupList []models.GroupChat
	if err := r.DB.Where("user_id = ?", input.UserID).Find(&groupList).Error; err != nil {
		return nil, err
	}
	return groupList, nil
}
