package domain

import (
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ID              int            `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID       *int           `json:"product_id"`
	Product         *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	Participants    []Participant  `json:"participants" gorm:"foreignKey:ConversationID"`
	Messages        []Message      `json:"messages" gorm:"foreignKey:ConversationID"`
	LastMessage     string         `json:"last_message,omitempty" gorm:"-"`
	LastMessageTime *time.Time     `json:"last_message_time,omitempty" gorm:"-"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
