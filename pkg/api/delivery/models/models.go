package models

type PrivateChat struct {
	UserID      string `json:"UserID"`
	RecipientID string `json:"RecipientID"`
}