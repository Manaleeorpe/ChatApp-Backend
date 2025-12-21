package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Messages struct {
	gorm.Model
	//ID uint `gorm:"primaryKey"` ID already exists in gorm.Model
	SenderID   uint
	Sender     User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;"`
	RecieverID uint
	Reciever   User `gorm:"foreignKey:RecieverID;constraint:OnDelete:CASCADE;"`
	Content    string
}

type MessagesHistory struct {
	gorm.Model
	//ID uint gorm:"primaryKey" ID already exists in gorm.Model
	ID           uint      `json:"ID"`
	Content      string    `json:"Content"`
	SenderID     uint      `json:"SenderID"`
	RecieverID   uint      `json:"RecieverID"`
	CreatedAt    time.Time `json:"CreatedAt"`
	SenderName   string    `json:"SenderName"`
	RecieverName string    `json:"RecieverName"`
}

func (M *Messages) SendAMessageModel() (*Messages, error) { //method is of struct type Message
	if err := db.Create(&M).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("Sender").Preload("Reciever").First(&M, M.ID).Error; err != nil {
		return nil, err
	}
	return M, nil
}

func GetMessageHistoryModel(ID1 int64, ID2 int64) ([]MessagesHistory, *gorm.DB) {
	var AllMessages []Messages

	dbResult := db.Preload("Sender").Preload("Reciever").
		Where("(sender_id = ? AND reciever_id = ?) OR (sender_id = ? AND reciever_id = ?)", ID1, ID2, ID2, ID1).
		Find(&AllMessages)

	//decide which to select friend1 or friend2 and get the required response
	var AllMessagesHistory []MessagesHistory

	for _, message := range AllMessages {
		r := MessagesHistory{
			ID:           message.ID,
			Content:      message.Content,
			SenderID:     message.SenderID,
			RecieverID:   message.RecieverID,
			CreatedAt:    message.CreatedAt,
			SenderName:   message.Sender.Name,
			RecieverName: message.Reciever.Name,
		}
		AllMessagesHistory = append(AllMessagesHistory, r)
	}

	//return response, dbResult
	return AllMessagesHistory, dbResult

}
