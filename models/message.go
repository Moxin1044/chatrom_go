package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type Message struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Username    string    `gorm:"size:100;not null" json:"username"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	MessageType string    `gorm:"size:50;not null" json:"message_type"` // text, image, file
	CreatedAt   time.Time `json:"timestamp"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移，创建表
	DB.AutoMigrate(&Message{})
}
