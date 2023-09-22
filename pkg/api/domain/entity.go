package domain

import (
	"time"

	"gorm.io/gorm"
)

type PrivateChat struct {
	gorm.Model
	UserID      string
	RecipientID string
	StartAt     time.Time
	LastSeen    time.Time
}

type PrivateChatHistory struct {
	gorm.Model
	UserID      string
	RecipientID string
	Text        string
	Status      string
	Time        time.Time
}

// Group chat

type GroupChat struct {
	gorm.Model
	UserID   string
	GroupID  string
	StartAt  time.Time
	LastSeen time.Time
}

type GroupChatHistory struct {
	gorm.Model
	UserID  string
	GroupID string
	Text    string
	Status  string
	Time    time.Time
}
