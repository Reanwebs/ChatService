package models

import (
	"time"
)

type PrivateChat struct {
	UserID      string `json:"UserID"`
	RecipientID string `json:"RecipientID"`
}

type GetChat struct {
	UserID string `json:"UserID"`
}

type PrivateChatHistory struct {
	Text   string `json:"Text"`
	Status string `json:"Status"`
	Time   time.Time
}
type PrivateChatWithHistory struct {
	PrivateChat
	PrivateChatHistory
}
