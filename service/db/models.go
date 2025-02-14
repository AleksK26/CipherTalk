package db

import (
	"time"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	PhotoURL string
}

type Conversation struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Members []User `gorm:"many2many:conversation_members;"`
}

type Message struct {
	ID             string `gorm:"primaryKey"`
	ConversationID string
	SenderID       string
	Content        string
	Timestamp      time.Time
	// Reactions     []string `gorm:"type:text[]"` // Optional: For reactions
}
