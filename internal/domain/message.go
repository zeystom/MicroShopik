package domain

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID             int            `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID int            `json:"conversation_id" gorm:"not null"`
	SenderID       *int           `json:"sender_id"`
	OrderID        *int           `json:"order_id"`
	Text           string         `json:"text" gorm:"type:text"`
	IsSystem       bool           `json:"is_system" gorm:"default:false"`
	Conversation   Conversation   `json:"conversation" gorm:"foreignKey:ConversationID"`
	Sender         *User          `json:"sender" gorm:"foreignKey:SenderID"`
	Order          *Order         `json:"order" gorm:"foreignKey:OrderID"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
