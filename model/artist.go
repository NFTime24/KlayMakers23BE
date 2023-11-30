package model

import "time"

type ArtistResult struct {
	Id           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	ProfilePath  string    `json:"profile_path"`
	Introduction string    `json:"introduction"`
	Instagram    string    `json:"instagram"`
}

type ArtistNameAndProfileResult struct {
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}

type Artist struct {
	DefaultSetting
	Name         string `gorm:"primaryKey;type:varchar"`
	Address      string `gorm:"primaryKey;type:varchar"`
	ProfilePath  string `gorm:"type:varchar"`
	Introduction string
	Instagram    string `gorm:"type:varchar"`
}
