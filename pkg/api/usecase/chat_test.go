package usecase

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPrivateChatRepo struct {
	mock.Mock
}

type MockGroupChatRepo struct {
	mock.Mock
}

var (
	mockRepo      = new(MockPrivateChatRepo)
	uc            = NewPrivateChatUsecase(mockRepo)
	mockGroupRepo = new(MockGroupChatRepo)
	gu            = NewGroupChatUsecase(mockGroupRepo)
)

func (m *MockPrivateChatRepo) CreatePrivateChat(entity domain.PrivateChat) error {
	entity.StartAt = time.Time{}
	entity.LastSeen = time.Time{}
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockPrivateChatRepo) AddPrivateChatHistory(entity domain.PrivateChatHistory) error {
	entity.Time = time.Time{}
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockPrivateChatRepo) GetChatList(UserID string) ([]models.PrivateChat, []models.PrivateChat, error) {
	args := m.Called(UserID)
	if args.Get(2) != nil {
		return nil, nil, args.Error(2)
	}

	existingChatList := []models.PrivateChat{}
	newChatList := []models.PrivateChat{}
	return existingChatList, newChatList, nil
}

func (m *MockPrivateChatRepo) GetPrivateChatHistory(UserID string, RecipientID string, DataLimitDays time.Time) ([]domain.PrivateChatHistory, error) {
	DataLimitDays = time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)
	args := m.Called(UserID, RecipientID, DataLimitDays)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	response := []domain.PrivateChatHistory{}
	return response, nil
}

func (m *MockPrivateChatRepo) GetRecievedChatHistory(UserID string, RecipientID string, DataLimitDays time.Time) ([]domain.PrivateChatHistory, error) {
	DataLimitDays = time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)
	args := m.Called(UserID, RecipientID, DataLimitDays)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	response := []domain.PrivateChatHistory{}
	return response, nil
}

func TestStartChat(t *testing.T) {
	repoInput := domain.PrivateChat{
		UserID:            "UserID",
		UserName:          "User1",
		RecipientID:       "RecipientID",
		RecipientName:     "User2",
		RecipientAvatarID: "RecipientAvatarID",
		NewRecipient:      true,
		StartAt:           time.Time{},
		LastSeen:          time.Time{},
	}

	usecaseInput := models.PrivateChat{
		UserID:            "UserID",
		UserName:          "User1",
		RecipientID:       "RecipientID",
		RecipientName:     "User2",
		RecipientAvatarID: "RecipientAvatarID",
		NewRecipient:      false,
		StartAt:           time.Time{},
		LastSeen:          time.Time{},
	}
	mockRepo.On("CreatePrivateChat", repoInput).Return(nil)

	err := uc.StartChat(usecaseInput)

	mockRepo.AssertCalled(t, "CreatePrivateChat", repoInput)

	assert.NoError(t, err)
}

func TestPrivateChatList(t *testing.T) {

	Input := models.GetChat{
		UserID: "UserID",
	}

	existingChatList := []models.PrivateChat{}
	newChatList := []models.PrivateChat{}

	mockRepo.On("GetChatList", Input.UserID).Return(existingChatList, newChatList, nil)

	uc := NewPrivateChatUsecase(mockRepo)

	combinedChatList, err := uc.PrivateChatList(Input)

	mockRepo.AssertCalled(t, "GetChatList", Input.UserID)

	assert.NoError(t, err)

	assert.Len(t, combinedChatList, len(existingChatList)+len(newChatList))
}

func TestCreatePrivateChatHistory(t *testing.T) {

	input := models.PrivateChatHistory{
		UserName:    "User1",
		UserID:      "UserID",
		RecipientID: "RecipientID",
		Text:        "Hello",
		Status:      "delivered",
		Time:        time.Time{},
	}

	repoInput := domain.PrivateChatHistory{
		UserName:    "User1",
		UserID:      "UserID",
		RecipientID: "RecipientID",
		Text:        "Hello",
		Status:      "delivered",
		Time:        time.Time{},
	}

	mockRepo.On("AddPrivateChatHistory", repoInput).Return(nil)

	err := uc.CreatePrivateChatHistory(input)

	mockRepo.AssertCalled(t, "AddPrivateChatHistory", repoInput)

	assert.NoError(t, err)

}

func TestRetrivePrivateChatHistory(t *testing.T) {
	DataLimitDays := time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)
	UserID := "UserID"
	RecipientID := "RecipientID"
	repoResponse := []domain.PrivateChatHistory{}
	useResponse := []models.PrivateChatHistory{}

	mockRepo.On("GetPrivateChatHistory", UserID, RecipientID, DataLimitDays).Return(repoResponse, nil)

	res, err := uc.RetrivePrivateChatHistory(UserID, RecipientID)

	mockRepo.AssertCalled(t, "GetPrivateChatHistory", UserID, RecipientID, DataLimitDays)

	assert.NoError(t, err)
	assert.Equal(t, res, useResponse)
}

func TestRetriveRecievedChatHistory(t *testing.T) {
	DataLimitDays := time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)
	UserID := "UserID"
	RecipientID := "RecipientID"
	repoResponse := []domain.PrivateChatHistory{}
	useResponse := []models.PrivateChatHistory{}

	mockRepo.On("GetRecievedChatHistory", UserID, RecipientID, DataLimitDays).Return(repoResponse, nil)

	res, err := uc.RetriveRecievedChatHistory(UserID, RecipientID)

	mockRepo.AssertCalled(t, "GetRecievedChatHistory", UserID, RecipientID, DataLimitDays)

	assert.NoError(t, err)
	assert.Equal(t, res, useResponse)
}

// Group Chat Unit Tests and Repo Mocks

func (m *MockGroupChatRepo) CreateGroupChat(entity domain.GroupChat) error {
	entity.StartAt = time.Time{}
	entity.LastSeen = time.Time{}
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockGroupChatRepo) GetGroupList(entity domain.GroupChat) ([]domain.GroupChat, error) {
	args := m.Called(entity)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	response := []domain.GroupChat{}
	return response, nil
}
func (m *MockGroupChatRepo) AddGroupChatHistory(entity domain.GroupChatHistory) error {
	entity.Time = time.Time{}
	args := m.Called(entity)
	return args.Error(0)
}
func (m *MockGroupChatRepo) GroupChatHistory(entity domain.GroupChatHistory, DataLimitDays time.Time) ([]domain.GroupChatHistory, error) {
	DataLimitDays = time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)
	args := m.Called(entity, DataLimitDays)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	response := []domain.GroupChatHistory{}
	return response, nil
}

func TestGroupChatStart(t *testing.T) {

	input := models.GroupChat{
		UserID:        "UserId",
		UserName:      "UserName",
		GroupID:       "GroupID",
		GroupName:     "GroupName",
		GroupAvatarID: "AvatarID",
		Permission:    false,
		StartAt:       time.Time{},
		LastSeen:      time.Time{},
	}
	entity := domain.GroupChat{
		UserID:        "UserId",
		UserName:      "UserName",
		GroupID:       "GroupID",
		GroupName:     "GroupName",
		GroupAvatarID: "AvatarID",
		Permission:    false,
		StartAt:       time.Time{},
		LastSeen:      time.Time{},
	}

	mockGroupRepo.On("CreateGroupChat", entity).Return(nil)

	err := gu.GroupChatStart(input)

	mockGroupRepo.AssertCalled(t, "CreateGroupChat", entity)

	assert.NoError(t, err)

}
func TestGetGroupList(t *testing.T) {
	repoInput := domain.GroupChat{
		UserID: "UserID",
	}
	input := models.GroupChat{
		UserID: "UserID",
	}
	mockGroupRepo.On("GetGroupList", repoInput).Return([]models.GroupChat{}, nil)

	res, err := gu.GetGroupList(input.UserID)

	mockGroupRepo.AssertCalled(t, "GetGroupList", repoInput)

	assert.NoError(t, err)
	assert.Empty(t, res)

}
func TestAddGroupChatHistory(t *testing.T) {
	entity := domain.GroupChatHistory{
		UserID:    "UserID",
		UserName:  "UserName",
		GroupID:   "GroupID",
		GroupName: "GroupName",
		Text:      "Hello",
		Status:    "delivered",
		Time:      time.Time{},
	}
	input := models.GroupChatHistory{
		UserID:    "UserID",
		UserName:  "UserName",
		GroupID:   "GroupID",
		GroupName: "GroupName",
		Text:      "Hello",
		Status:    "delivered",
		Time:      time.Time{},
	}

	mockGroupRepo.On("AddGroupChatHistory", entity).Return(nil)

	err := gu.AddGroupChatHistory(input)

	mockGroupRepo.AssertCalled(t, "AddGroupChatHistory", entity)

	assert.NoError(t, err)
}
func TestGetGroupChatHistory(t *testing.T) {
	entity := domain.GroupChatHistory{
		GroupID: "GroupID",
	}
	input := models.GroupChatHistory{
		GroupID: "GroupID",
	}
	DataLimitDays := time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)

	mockGroupRepo.On("GroupChatHistory", entity, DataLimitDays).Return([]domain.GroupChatHistory{}, nil)

	res, err := gu.GetGroupChatHistory(input)

	mockGroupRepo.AssertCalled(t, "GroupChatHistory", entity, DataLimitDays)

	assert.Empty(t, res)
	assert.NoError(t, err)
}
