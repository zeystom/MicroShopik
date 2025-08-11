package domain

import (
	"gorm.io/gorm"
	"time"
)

type Participant struct {
	ConversationID int            `json:"conversation_id" gorm:"primaryKey"`
	UserID         int            `json:"user_id" gorm:"primaryKey"`
	Conversation   Conversation   `json:"conversation" gorm:"foreignKey:ConversationID"`
	User           User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
