package main

import (
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (m *message) createMessage(db *gorm.DB) error {
	err := db.Create(&m).Error
	if err != nil {
		return err
	}

	return nil
}

func getMessages(db *gorm.DB, start, count int) ([]message, error) {
	var messages []message
	var err error
	if start == 0 || count == 0 {
		err = db.Omit("contact").Order("id desc").Find(&messages).Error
	} else {
		err = db.Omit("contact").Order("id desc").Limit(count).Offset(start).Find(&messages).Error
	}

	if err != nil {
		return nil, err
	}

	return messages, nil
}
