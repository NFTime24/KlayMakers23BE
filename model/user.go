package model

import (
	"time"
)

type Log struct {
	DefaultSetting
	UserAddress string `gorm:"type:varchar"`
	Status      string `gorm:"type:varchar"`
}

type User struct {
	DefaultSetting
	UserAddress string `gorm:"primaryKey;type:varchar"`
	NickName    string `gorm:"type:varchar"`
	ProfilePath string `gorm:"type:varchar"`
}

type UserResult struct {
	Id          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserAddress string    `json:"user_address"`
	NickName    string    `json:"nick_name"`
	ProfilePath string    `json:"profile_path"`
	TicketCount uint      `json:"ticket_count"`
}
