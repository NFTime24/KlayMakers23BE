package model

import "time"

type DefaultSetting struct {
	ID        uint      `gorm:"autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp without time zone"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp without time zone"`
}
