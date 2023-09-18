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

type ChatHistory struct {
	gorm.Model
	UserID      string
	RecipientID string
}
