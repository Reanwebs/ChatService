package domain

import (
	"time"

	"gorm.io/gorm"
)

type PrivateChat struct {
	gorm.Model
	UserID            string
	UserName          string
	RecipientID       string
	RecipientName     string
	RecipientAvatarID string
	NewRecipient      bool `gorm:"default:true"`
	StartAt           time.Time
	LastSeen          time.Time
}

type PrivateChatHistory struct {
	gorm.Model
	UserName    string
	UserID      string
	RecipientID string
	Text        string
	Status      string
	Time        time.Time
}

// Group chat

type GroupChat struct {
	gorm.Model
	UserID        string
	UserName      string
	GroupID       string
	GroupName     string
	GroupAvatarID string
	StartAt       time.Time
	LastSeen      time.Time
}

type GroupChatHistory struct {
	gorm.Model
	UserID    string
	UserName  string
	GroupID   string
	GroupName string
	Text      string
	Status    string
	Time      time.Time
}
