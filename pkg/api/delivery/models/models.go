package models

import (
	"time"
)

type StartChat struct {
	UserName    string `json:"UserName"`
	RecipientID string `json:"RecipientID"`
}
type GetChat struct {
	UserID string `json:"UserID"`
}

type ChatHistory struct {
	RecipientID string `json:"RecipientID"`
}
type PrivateChat struct {
	UserID            string `json:"UserID"`
	UserName          string `json:"UserName"`
	RecipientID       string `json:"RecipientID"`
	RecipientName     string `json:"RecipientName"`
	RecipientAvatarID string `json:"AvatarID"`
	NewRecipient      bool   `json:"-"`
	StartAt           time.Time
	LastSeen          time.Time
}

type PrivateChatHistory struct {
	UserName    string `json:"UserName"`
	UserID      string `json:"UserID"`
	RecipientID string `json:"RecipientID"`
	Text        string `json:"Text"`
	Status      string `json:"Status"`
	Time        time.Time
}

type ChatHistoryResponse struct {
	Messages []PrivateChatHistory `json:"messages"`
}

// Group

type GroupChat struct {
	UserID        string `json:"UserID"`
	UserName      string `jaon:"UserName"`
	GroupID       string `json:"GroupID"`
	GroupName     string `json:"GroupName"`
	GroupAvatarID string `json:"AvaterID"`
	StartAt       time.Time
	LastSeen      time.Time
}

type GetGroupChat struct {
	UserID string `json:"UserID"`
}

type GroupChatHistory struct {
	UserID    string `json:"UserID"`
	UserName  string `jaon:"UserName"`
	GroupID   string `json:"GroupID"`
	GroupName string `json:"GroupName"`
	Text      string `json:"Text"`
	Status    string `json:"Status"`
	Time      time.Time
}
