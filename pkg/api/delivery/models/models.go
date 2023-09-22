package models

import (
	"time"
)

type PrivateChat struct {
	UserID      string `json:"UserID"`
	RecipientID string `json:"RecipientID"`
	StartAt     time.Time
	LastSeen    time.Time
}

type GetChat struct {
	UserID string `json:"UserID"`
}

type PrivateChatHistory struct {
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
	UserID   string `json:"UserID"`
	GroupID  string `json:"GroupID"`
	StartAt  time.Time
	LastSeen time.Time
}

type GetGroupChat struct {
	UserID string `json:"UserID"`
}

type GroupChatHistory struct {
	GroupID string `json:"GroupID"`
	Text    string `json:"Text"`
	Status  string `json:"Status"`
	Time    time.Time
}
