package domain

import (
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Participants []Participant  `json:"participants" gorm:"foreignKey:ConversationID"`
	Messages     []Message      `json:"messages" gorm:"foreignKey:ConversationID"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
